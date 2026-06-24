package iris

import (
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sylunbelievable/secondadmin/server/internal/config"
	"github.com/sylunbelievable/secondadmin/server/internal/dto"
	"github.com/sylunbelievable/secondadmin/server/internal/service"
	"go.uber.org/zap"
)

func registerRoutes(app *iris.Application, services *service.Services, cfg config.Config, log *zap.Logger) {
	api := app.Party("/api/v1")
	auth := api.Party("/auth")
	auth.Post("/login", login(services, cfg))
	auth.Post("/refresh", refresh(services, cfg))

	account := auth.Party("", authenticate(services), audit(services, log))
	account.Post("/logout", logout(services, cfg))
	account.Get("/me", me(services))
	account.Get("/sessions", sessions(services))
	account.Delete("/sessions/{id:string}", deleteSession(services))

	protected := api.Party("", authenticate(services))
	protected.Get("/menus/current", currentMenus(services))
	protected.Get("/dictionaries/{code:string}/items", dictionaryItems(services))

	admin := api.Party("", authenticate(services), audit(services, log), authorize(services))
	admin.Get("/users", listUsers(services))
	admin.Post("/users", createUser(services))
	admin.Get("/users/{id:uint64}", getUser(services))
	admin.Put("/users/{id:uint64}", updateUser(services))
	admin.Delete("/users/{id:uint64}", disableUser(services))
	admin.Put("/users/{id:uint64}/roles", setUserRoles(services))

	admin.Get("/roles", listRoles(services))
	admin.Post("/roles", createRole(services))
	admin.Put("/roles/{id:uint64}", updateRole(services))
	admin.Delete("/roles/{id:uint64}", deleteRole(services))
	admin.Put("/roles/{id:uint64}/apis", setRoleAPIs(services))

	admin.Get("/apis", listAPIs(services))
	admin.Post("/apis", createAPI(services))
	admin.Put("/apis/{id:uint64}", updateAPI(services))
	admin.Delete("/apis/{id:uint64}", deleteAPI(services))

	registerSystemRoutes(admin, services)
}

func login(s *service.Services, cfg config.Config) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.LoginRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		tokens, user, err := s.Login(ctx.Request().Context(), req, ctx.RemoteAddr(), ctx.GetHeader("User-Agent"))
		if err != nil {
			writeError(ctx, err)
			return
		}
		if req.AuthMode == "cookie" {
			setAuthCookies(ctx, cfg, tokens)
			tokens.AccessToken, tokens.RefreshToken = "", ""
		}
		_ = ctx.JSON(iris.Map{"tokens": tokens, "user": user})
	}
}

func refresh(s *service.Services, cfg config.Config) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.RefreshRequest
		_ = ctx.ReadJSON(&req)
		cookieMode := req.RefreshToken == ""
		if cookieMode {
			req.RefreshToken = ctx.GetCookie("refresh_token")
			sessionID, _, ok := stringsCut(req.RefreshToken)
			if !ok {
				writeError(ctx, service.ErrForbidden)
				return
			}
			if err := s.CheckCSRF(ctx.Request().Context(), sessionID, ctx.GetHeader("X-CSRF-Token")); err != nil {
				writeError(ctx, err)
				return
			}
		}
		tokens, p, err := s.Refresh(ctx.Request().Context(), req.RefreshToken, ctx.RemoteAddr(), ctx.GetHeader("User-Agent"))
		if err != nil {
			writeError(ctx, err)
			return
		}
		if cookieMode || p.AuthMode == "cookie" {
			setAuthCookies(ctx, cfg, tokens)
			tokens.AccessToken, tokens.RefreshToken = "", ""
		}
		_ = ctx.JSON(tokens)
	}
}

func logout(s *service.Services, cfg config.Config) iris.Handler {
	return func(ctx iris.Context) {
		if err := s.Logout(ctx.Request().Context(), principal(ctx), ctx.RemoteAddr(), ctx.GetHeader("User-Agent")); err != nil {
			writeError(ctx, err)
			return
		}
		clearAuthCookies(ctx, cfg)
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func me(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		user, err := s.User(ctx.Request().Context(), principal(ctx).UserID)
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(user)
	}
}

func sessions(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		items, err := s.Sessions(ctx.Request().Context(), principal(ctx).UserID)
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(items)
	}
}

func deleteSession(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		if err := s.DeleteSession(ctx.Request().Context(), principal(ctx).UserID, ctx.Params().Get("id")); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func listUsers(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		page, size := positive(ctx.URLParamIntDefault("page", 1), 1), positive(ctx.URLParamIntDefault("pageSize", 20), 20)
		if size > 100 {
			size = 100
		}
		items, total, err := s.ListUsers(ctx.Request().Context(), page, size)
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(iris.Map{"items": items, "total": total})
	}
}

func createUser(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.CreateUserRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		item, err := s.CreateUser(ctx.Request().Context(), req)
		if err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(item)
	}
}

func getUser(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		item, err := s.User(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0))
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(item)
	}
}

func updateUser(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.UpdateUserRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.UpdateUser(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), req); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func disableUser(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		status := int16(0)
		if err := s.UpdateUser(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), dto.UpdateUserRequest{Status: &status}); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func setUserRoles(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.IDsRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.SetUserRoles(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), []uint64(req.IDs)); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func listRoles(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		items, err := s.ListRoles(ctx.Request().Context())
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(items)
	}
}

func createRole(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.CreateRoleRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		item, err := s.CreateRole(ctx.Request().Context(), req)
		if err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(item)
	}
}

func updateRole(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.UpdateRoleRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.UpdateRole(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), req); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func deleteRole(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		if err := s.DeleteRole(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0)); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func setRoleAPIs(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.IDsRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.SetRoleAPIs(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), []uint64(req.IDs)); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func listAPIs(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		items, err := s.ListAPIs(ctx.Request().Context())
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(items)
	}
}

func createAPI(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.CreateAPIRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		item, err := s.CreateAPI(ctx.Request().Context(), req)
		if err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(item)
	}
}

func updateAPI(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.CreateAPIRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.UpdateAPI(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), req); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func deleteAPI(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		if err := s.DeleteAPI(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0)); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func setAuthCookies(ctx iris.Context, cfg config.Config, tokens dto.Tokens) {
	sameSite := http.SameSiteLaxMode
	if cfg.Auth.Cookie.SameSite == "strict" {
		sameSite = http.SameSiteStrictMode
	}
	http.SetCookie(ctx.ResponseWriter(), &http.Cookie{Name: "access_token", Value: tokens.AccessToken, Path: "/", HttpOnly: true, Secure: cfg.Auth.Cookie.Secure, SameSite: sameSite, MaxAge: int(cfg.Auth.AccessTokenTTL.Seconds())})
	http.SetCookie(ctx.ResponseWriter(), &http.Cookie{Name: "refresh_token", Value: tokens.RefreshToken, Path: "/api/v1/auth/refresh", HttpOnly: true, Secure: cfg.Auth.Cookie.Secure, SameSite: sameSite, MaxAge: int(cfg.Auth.RefreshTokenTTL.Seconds())})
	http.SetCookie(ctx.ResponseWriter(), &http.Cookie{Name: "csrf_token", Value: tokens.CSRFToken, Path: "/", Secure: cfg.Auth.Cookie.Secure, SameSite: sameSite, MaxAge: int(cfg.Auth.RefreshTokenTTL.Seconds())})
}

func clearAuthCookies(ctx iris.Context, cfg config.Config) {
	for _, cookie := range []http.Cookie{
		{Name: "access_token", Path: "/"}, {Name: "refresh_token", Path: "/api/v1/auth/refresh"}, {Name: "csrf_token", Path: "/"},
	} {
		cookie.MaxAge, cookie.Expires, cookie.HttpOnly, cookie.Secure = -1, time.Unix(1, 0), cookie.Name != "csrf_token", cfg.Auth.Cookie.Secure
		http.SetCookie(ctx.ResponseWriter(), &cookie)
	}
}

func positive(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

func stringsCut(value string) (string, string, bool) {
	for i := range value {
		if value[i] == '.' {
			return value[:i], value[i+1:], true
		}
	}
	return "", "", false
}
