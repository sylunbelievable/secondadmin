package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sylunbelievable/secondadmin/server/internal/config"
	"github.com/sylunbelievable/secondadmin/server/internal/dto"
	"github.com/sylunbelievable/secondadmin/server/internal/entity"
	"github.com/sylunbelievable/secondadmin/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidInput       = errors.New("invalid input")
	ErrConflict           = repository.ErrConflict
	ErrNotFound           = repository.ErrNotFound
	ErrDependency         = errors.New("dependency unavailable")
	ErrRateLimited        = errors.New("too many login attempts")
)

type Services struct {
	Repos  *repository.Repositories
	Redis  *redis.Client
	Casbin *casbin.Enforcer
	Config config.Config
}

type Principal struct {
	UserID    uint64
	SessionID string
	AuthMode  string
}

type accessClaims struct {
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

func New(repos *repository.Repositories, redisClient *redis.Client, enforcer *casbin.Enforcer, cfg config.Config) *Services {
	return &Services{Repos: repos, Redis: redisClient, Casbin: enforcer, Config: cfg}
}

func HashPassword(password string) (string, error) {
	if len(password) < 8 || len(password) > 72 {
		return "", fmt.Errorf("%w: password must be 8-72 bytes", ErrInvalidInput)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *Services) Login(ctx context.Context, req dto.LoginRequest, ip, userAgent string) (dto.Tokens, dto.User, error) {
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	if req.AuthMode != "cookie" && req.AuthMode != "bearer" || req.Username == "" || req.Password == "" {
		return dto.Tokens{}, dto.User{}, ErrInvalidInput
	}
	failureKey := "auth:login:" + ip + ":" + req.Username
	// ponytail: fixed 5/5m limit; add config only when deployments need different thresholds.
	failures, failureErr := s.Redis.Get(ctx, failureKey).Int()
	if failureErr != nil && !errors.Is(failureErr, redis.Nil) {
		return dto.Tokens{}, dto.User{}, fmt.Errorf("%w: redis: %v", ErrDependency, failureErr)
	}
	if failures >= 5 {
		return dto.Tokens{}, dto.User{}, ErrRateLimited
	}
	user, err := s.Repos.UserByUsername(ctx, req.Username)
	if err != nil || user.Status != 1 || !CheckPassword(user.PasswordHash, req.Password) {
		if count, incrementErr := s.Redis.Incr(ctx, failureKey).Result(); incrementErr == nil && count == 1 {
			_ = s.Redis.Expire(ctx, failureKey, 5*time.Minute).Err()
		}
		s.loginLog(ctx, nil, req.Username, "login", false, ip, userAgent, req.DeviceID)
		return dto.Tokens{}, dto.User{}, ErrInvalidCredentials
	}
	_ = s.Redis.Del(ctx, failureKey).Err()

	sessionID := uuid.NewString()
	refreshSecret, err := randomToken()
	if err != nil {
		return dto.Tokens{}, dto.User{}, err
	}
	refresh := sessionID + "." + refreshSecret
	csrf, err := randomToken()
	if err != nil {
		return dto.Tokens{}, dto.User{}, err
	}
	if req.DeviceID == "" {
		req.DeviceID = uuid.NewString()
	}
	if err := s.createSession(ctx, user.ID, sessionID, req.DeviceID, req.AuthMode, refresh, csrf); err != nil {
		return dto.Tokens{}, dto.User{}, fmt.Errorf("%w: redis: %v", ErrDependency, err)
	}
	access, err := s.signAccessToken(user.ID, sessionID)
	if err != nil {
		_ = s.deleteSession(ctx, user.ID, sessionID)
		return dto.Tokens{}, dto.User{}, err
	}
	s.loginLog(ctx, &user.ID, user.Username, "login", true, ip, userAgent, req.DeviceID)
	return dto.Tokens{
		AccessToken: access, RefreshToken: refresh, CSRFToken: csrf,
		ExpiresIn: int64(s.Config.Auth.AccessTokenTTL.Seconds()),
	}, toUserDTO(user), nil
}

func (s *Services) Refresh(ctx context.Context, refreshToken, ip, userAgent string) (dto.Tokens, Principal, error) {
	if refreshToken == "" {
		return dto.Tokens{}, Principal{}, ErrUnauthorized
	}
	sessionID, _, ok := strings.Cut(refreshToken, ".")
	if !ok || sessionID == "" {
		return dto.Tokens{}, Principal{}, ErrUnauthorized
	}
	key := sessionKey(sessionID)
	values, err := s.Redis.HMGet(ctx, key, "user_id", "auth_mode").Result()
	if err != nil {
		return dto.Tokens{}, Principal{}, fmt.Errorf("%w: redis: %v", ErrDependency, err)
	}
	if values[0] == nil {
		return dto.Tokens{}, Principal{}, ErrUnauthorized
	}
	userID, _ := strconv.ParseUint(values[0].(string), 10, 64)
	user, err := s.Repos.UserByID(ctx, userID)
	if err != nil || user.Status != 1 {
		_ = s.deleteSession(ctx, userID, sessionID)
		return dto.Tokens{}, Principal{}, ErrUnauthorized
	}
	nextSecret, err := randomToken()
	if err != nil {
		return dto.Tokens{}, Principal{}, err
	}
	nextRefresh := sessionID + "." + nextSecret
	csrf, err := randomToken()
	if err != nil {
		return dto.Tokens{}, Principal{}, err
	}
	rotated, err := rotateSessionScript.Run(ctx, s.Redis, []string{key},
		tokenHash(refreshToken), tokenHash(nextRefresh), tokenHash(csrf),
		int64(s.Config.Auth.RefreshTokenTTL.Seconds())).Int()
	if err != nil {
		return dto.Tokens{}, Principal{}, fmt.Errorf("%w: redis: %v", ErrDependency, err)
	}
	if rotated != 1 {
		return dto.Tokens{}, Principal{}, ErrUnauthorized
	}
	access, err := s.signAccessToken(userID, sessionID)
	if err != nil {
		return dto.Tokens{}, Principal{}, err
	}
	mode, _ := values[1].(string)
	s.loginLog(ctx, &userID, user.Username, "refresh", true, ip, userAgent, "")
	return dto.Tokens{AccessToken: access, RefreshToken: nextRefresh, CSRFToken: csrf, ExpiresIn: int64(s.Config.Auth.AccessTokenTTL.Seconds())},
		Principal{UserID: userID, SessionID: sessionID, AuthMode: mode}, nil
}

func (s *Services) Authenticate(ctx context.Context, token string) (Principal, error) {
	claims, err := s.parseAccessToken(token)
	if err != nil || claims.Subject == "" || claims.SessionID == "" {
		return Principal{}, ErrUnauthorized
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return Principal{}, ErrUnauthorized
	}
	mode, err := s.Redis.HGet(ctx, sessionKey(claims.SessionID), "auth_mode").Result()
	if errors.Is(err, redis.Nil) {
		return Principal{}, ErrUnauthorized
	}
	if err != nil {
		return Principal{}, fmt.Errorf("%w: redis: %v", ErrDependency, err)
	}
	return Principal{UserID: userID, SessionID: claims.SessionID, AuthMode: mode}, nil
}

func (s *Services) CheckCSRF(ctx context.Context, sessionID, token string) error {
	stored, err := s.Redis.HGet(ctx, sessionKey(sessionID), "csrf_hash").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("%w: redis: %v", ErrDependency, err)
	}
	if err != nil || tokenHash(token) != stored {
		return ErrForbidden
	}
	return nil
}

func (s *Services) Logout(ctx context.Context, p Principal, ip, userAgent string) error {
	if err := s.deleteSession(ctx, p.UserID, p.SessionID); err != nil {
		return fmt.Errorf("%w: redis: %v", ErrDependency, err)
	}
	if user, err := s.Repos.UserByID(ctx, p.UserID); err == nil {
		s.loginLog(ctx, &p.UserID, user.Username, "logout", true, ip, userAgent, "")
	}
	return nil
}

func (s *Services) Sessions(ctx context.Context, userID uint64) ([]dto.Session, error) {
	ids, err := s.Redis.ZRange(ctx, userSessionsKey(userID), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("%w: redis: %v", ErrDependency, err)
	}
	sessions := make([]dto.Session, 0, len(ids))
	for _, id := range ids {
		values, err := s.Redis.HMGet(ctx, sessionKey(id), "device_id", "auth_mode", "created_at").Result()
		if err != nil || values[0] == nil {
			continue
		}
		created, _ := time.Parse(time.RFC3339Nano, values[2].(string))
		sessions = append(sessions, dto.Session{ID: id, DeviceID: values[0].(string), AuthMode: values[1].(string), CreatedAt: created})
	}
	return sessions, nil
}

func (s *Services) DeleteSession(ctx context.Context, userID uint64, sessionID string) error {
	if _, err := s.Redis.ZScore(ctx, userSessionsKey(userID), sessionID).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			return fmt.Errorf("%w: redis: %v", ErrDependency, err)
		}
		return ErrForbidden
	}
	if err := s.deleteSession(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("%w: redis: %v", ErrDependency, err)
	}
	return nil
}

func (s *Services) User(ctx context.Context, id uint64) (dto.User, error) {
	user, err := s.Repos.UserByID(ctx, id)
	return toUserDTO(user), err
}

func (s *Services) ListUsers(ctx context.Context, page, size int) ([]dto.User, int64, error) {
	users, total, err := s.Repos.ListUsers(ctx, page, size)
	out := make([]dto.User, len(users))
	for i := range users {
		out[i] = toUserDTO(users[i])
	}
	return out, total, err
}

func (s *Services) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.User, error) {
	hash, err := HashPassword(req.Password)
	if err != nil {
		return dto.User{}, err
	}
	user := entity.User{Username: req.Username, PasswordHash: hash, Nickname: strings.TrimSpace(req.Nickname), Status: 1, PasswordChangedAt: time.Now()}
	if err := s.Repos.CreateUser(ctx, &user); err != nil {
		return dto.User{}, err
	}
	return toUserDTO(user), nil
}

func (s *Services) UpdateUser(ctx context.Context, id uint64, req dto.UpdateUserRequest) error {
	values := map[string]any{}
	if req.Nickname != nil {
		values["nickname"] = strings.TrimSpace(*req.Nickname)
	}
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return ErrInvalidInput
		}
		values["status"] = *req.Status
	}
	if req.Password != nil {
		hash, err := HashPassword(*req.Password)
		if err != nil {
			return err
		}
		values["password_hash"], values["password_changed_at"] = hash, time.Now()
	}
	if len(values) == 0 {
		return ErrInvalidInput
	}
	if err := s.Repos.UpdateUser(ctx, id, values); err != nil {
		return err
	}
	if req.Status != nil && *req.Status == 0 {
		if err := s.deleteAllSessions(ctx, id); err != nil {
			return fmt.Errorf("%w: redis: %v", ErrDependency, err)
		}
	}
	return nil
}

func (s *Services) SetUserRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	if _, err := s.Repos.UserByID(ctx, userID); err != nil {
		return err
	}
	roles := make([]string, 0, len(roleIDs))
	for _, id := range roleIDs {
		role, err := s.Repos.RoleByID(ctx, id)
		if err != nil || role.Status != 1 {
			return ErrInvalidInput
		}
		roles = append(roles, role.Code)
	}
	subject := strconv.FormatUint(userID, 10)
	if _, err := s.Casbin.DeleteRolesForUser(subject); err != nil {
		return err
	}
	for _, role := range roles {
		if _, err := s.Casbin.AddRoleForUser(subject, role); err != nil {
			return err
		}
	}
	return nil
}

func (s *Services) CreateRole(ctx context.Context, req dto.CreateRoleRequest) (dto.Role, error) {
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
		return dto.Role{}, ErrInvalidInput
	}
	role := entity.Role{Code: req.Code, Name: strings.TrimSpace(req.Name), Status: 1}
	if err := s.Repos.CreateRole(ctx, &role); err != nil {
		return dto.Role{}, err
	}
	return toRoleDTO(role), nil
}

func (s *Services) ListRoles(ctx context.Context) ([]dto.Role, error) {
	roles, err := s.Repos.ListRoles(ctx)
	out := make([]dto.Role, len(roles))
	for i := range roles {
		out[i] = toRoleDTO(roles[i])
	}
	return out, err
}

func (s *Services) UpdateRole(ctx context.Context, id uint64, req dto.UpdateRoleRequest) error {
	role, err := s.Repos.RoleByID(ctx, id)
	if err != nil {
		return err
	}
	if role.Code == superAdminRole && req.Status != nil && *req.Status == 0 {
		return ErrConflict
	}
	values := map[string]any{}
	if req.Name != nil {
		values["name"] = strings.TrimSpace(*req.Name)
	}
	if req.Status != nil {
		values["status"] = *req.Status
	}
	if len(values) == 0 {
		return ErrInvalidInput
	}
	if err := s.Repos.UpdateRole(ctx, id, values); err != nil {
		return err
	}
	if req.Status != nil && *req.Status == 0 {
		_, err = s.Casbin.DeleteRole(role.Code)
		return err
	}
	return nil
}

func (s *Services) DeleteRole(ctx context.Context, id uint64) error {
	role, err := s.Repos.RoleByID(ctx, id)
	if err != nil {
		return err
	}
	if role.Code == superAdminRole {
		return ErrConflict
	}
	if _, err := s.Casbin.DeleteRole(role.Code); err != nil {
		return err
	}
	return s.Repos.DeleteRole(ctx, id)
}

func (s *Services) SetRoleAPIs(ctx context.Context, roleID uint64, apiIDs []uint64) error {
	role, err := s.Repos.RoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if _, err := s.Casbin.RemoveFilteredPolicy(0, role.Code); err != nil {
		return err
	}
	for _, id := range apiIDs {
		api, err := s.Repos.APIByID(ctx, id)
		if err != nil {
			return ErrInvalidInput
		}
		if _, err := s.Casbin.AddPolicy(role.Code, api.Path, api.Method); err != nil {
			return err
		}
	}
	return nil
}

func (s *Services) CreateAPI(ctx context.Context, req dto.CreateAPIRequest) (dto.API, error) {
	if req.Path == "" || req.Method == "" || req.Name == "" {
		return dto.API{}, ErrInvalidInput
	}
	api := entity.API{Group: strings.TrimSpace(req.Group), Name: strings.TrimSpace(req.Name), Path: req.Path, Method: req.Method}
	if err := s.Repos.CreateAPI(ctx, &api); err != nil {
		return dto.API{}, err
	}
	return toAPIDTO(api), nil
}

func (s *Services) ListAPIs(ctx context.Context) ([]dto.API, error) {
	apis, err := s.Repos.ListAPIs(ctx)
	out := make([]dto.API, len(apis))
	for i := range apis {
		out[i] = toAPIDTO(apis[i])
	}
	return out, err
}

func (s *Services) UpdateAPI(ctx context.Context, id uint64, req dto.CreateAPIRequest) error {
	if req.Path == "" || req.Method == "" || req.Name == "" {
		return ErrInvalidInput
	}
	old, err := s.Repos.APIByID(ctx, id)
	if err != nil {
		return err
	}
	roles, err := s.Casbin.GetFilteredPolicy(1, old.Path, old.Method)
	if err != nil {
		return err
	}
	if err := s.Repos.UpdateAPI(ctx, id, map[string]any{
		"group": strings.TrimSpace(req.Group), "name": strings.TrimSpace(req.Name),
		"path": req.Path, "method": strings.ToUpper(req.Method),
	}); err != nil {
		return err
	}
	if _, err := s.Casbin.RemoveFilteredPolicy(1, old.Path, old.Method); err != nil {
		return err
	}
	for _, policy := range roles {
		if len(policy) > 0 {
			if _, err := s.Casbin.AddPolicy(policy[0], req.Path, strings.ToUpper(req.Method)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Services) DeleteAPI(ctx context.Context, id uint64) error {
	api, err := s.Repos.APIByID(ctx, id)
	if err != nil {
		return err
	}
	if _, err := s.Casbin.RemoveFilteredPolicy(1, api.Path, api.Method); err != nil {
		return err
	}
	return s.Repos.DeleteAPI(ctx, id)
}

func (s *Services) Authorize(userID uint64, path, method string) error {
	subject := strconv.FormatUint(userID, 10)
	roles, err := s.Casbin.GetRolesForUser(subject)
	if err != nil {
		return err
	}
	if slices.Contains(roles, superAdminRole) {
		return nil
	}
	ok, err := s.Casbin.Enforce(subject, path, method)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

func (s *Services) signAccessToken(userID uint64, sessionID string) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims{
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.Config.Auth.AccessTokenTTL)),
		},
	}).SignedString([]byte(s.Config.Auth.JWTSecret))
}

func (s *Services) parseAccessToken(token string) (accessClaims, error) {
	claims := accessClaims{}
	parsed, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnauthorized
		}
		return []byte(s.Config.Auth.JWTSecret), nil
	})
	if err != nil || !parsed.Valid {
		return accessClaims{}, ErrUnauthorized
	}
	return claims, nil
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func tokenHash(token string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(token))) }
func sessionKey(id string) string   { return "auth:session:" + id }
func userSessionsKey(id uint64) string {
	return "auth:user:" + strconv.FormatUint(id, 10) + ":sessions"
}

var createSessionScript = redis.NewScript(`
redis.call("HSET", KEYS[1], "user_id", ARGV[1], "refresh_hash", ARGV[2], "csrf_hash", ARGV[3], "device_id", ARGV[4], "auth_mode", ARGV[5], "created_at", ARGV[6])
redis.call("EXPIRE", KEYS[1], ARGV[7])
redis.call("ZADD", KEYS[2], ARGV[8], ARGV[9])
redis.call("EXPIRE", KEYS[2], ARGV[7])
local max = tonumber(ARGV[10])
if max > 0 then
  local excess = redis.call("ZCARD", KEYS[2]) - max
  if excess > 0 then
    local old = redis.call("ZRANGE", KEYS[2], 0, excess - 1)
    for _, id in ipairs(old) do redis.call("DEL", "auth:session:" .. id) end
    redis.call("ZREM", KEYS[2], unpack(old))
  end
end
return 1
`)

var rotateSessionScript = redis.NewScript(`
if redis.call("HGET", KEYS[1], "refresh_hash") ~= ARGV[1] then return 0 end
redis.call("HSET", KEYS[1], "refresh_hash", ARGV[2], "csrf_hash", ARGV[3])
redis.call("EXPIRE", KEYS[1], ARGV[4])
return 1
`)

func (s *Services) createSession(ctx context.Context, userID uint64, sessionID, deviceID, mode, refresh, csrf string) error {
	now := time.Now()
	return createSessionScript.Run(ctx, s.Redis, []string{sessionKey(sessionID), userSessionsKey(userID)},
		userID, tokenHash(refresh), tokenHash(csrf), deviceID, mode, now.Format(time.RFC3339Nano),
		int64(s.Config.Auth.RefreshTokenTTL.Seconds()), now.UnixNano(), sessionID, s.Config.Auth.MaxDevices).Err()
}

func (s *Services) deleteSession(ctx context.Context, userID uint64, sessionID string) error {
	_, err := s.Redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, sessionKey(sessionID))
		pipe.ZRem(ctx, userSessionsKey(userID), sessionID)
		return nil
	})
	return err
}

func (s *Services) deleteAllSessions(ctx context.Context, userID uint64) error {
	ids, err := s.Redis.ZRange(ctx, userSessionsKey(userID), 0, -1).Result()
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(ids)+1)
	for _, id := range ids {
		keys = append(keys, sessionKey(id))
	}
	keys = append(keys, userSessionsKey(userID))
	return s.Redis.Del(ctx, keys...).Err()
}

func (s *Services) loginLog(ctx context.Context, userID *uint64, username, event string, success bool, ip, userAgent, deviceID string) {
	_ = s.Repos.AddLoginLog(ctx, &entity.LoginLog{
		UserID: userID, Username: username, Event: event, Success: success,
		IP: ip, UserAgent: userAgent, DeviceID: deviceID,
	})
}

func toUserDTO(user entity.User) dto.User {
	return dto.User{ID: user.ID, Username: user.Username, Nickname: user.Nickname, Status: user.Status}
}

func toRoleDTO(role entity.Role) dto.Role {
	return dto.Role{ID: role.ID, Code: role.Code, Name: role.Name, Status: role.Status}
}

func toAPIDTO(api entity.API) dto.API {
	return dto.API{ID: api.ID, Group: api.Group, Name: api.Name, Path: api.Path, Method: api.Method}
}
