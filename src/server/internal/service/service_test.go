package service

import (
	"errors"
	"testing"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/sylunbelievable/secondadmin/server/internal/config"
)

func TestPasswordHash(t *testing.T) {
	hash, err := HashPassword("correct-horse")
	if err != nil || !CheckPassword(hash, "correct-horse") || CheckPassword(hash, "wrong") {
		t.Fatal("password hashing check failed")
	}
}

func TestAccessToken(t *testing.T) {
	s := Services{Config: config.Config{Auth: config.Auth{JWTSecret: "test-secret", AccessTokenTTL: time.Minute}}}
	token, err := s.signAccessToken(42, "session")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := s.parseAccessToken(token)
	if err != nil || claims.Subject != "42" || claims.SessionID != "session" {
		t.Fatal("access token round trip failed")
	}
}

func TestRefreshTokenRotationMaterial(t *testing.T) {
	a, err := randomToken()
	if err != nil {
		t.Fatal(err)
	}
	b, err := randomToken()
	if err != nil || a == b || tokenHash(a) == tokenHash(b) {
		t.Fatal("refresh token material must rotate")
	}
}

func TestAuthorizeAllowsSuperAdminWithoutPolicy(t *testing.T) {
	enforcer := testEnforcer(t)
	if _, err := enforcer.AddRoleForUser("1", superAdminRole); err != nil {
		t.Fatal(err)
	}
	s := Services{Casbin: enforcer}
	if err := s.Authorize(1, "/api/v1/anything", "DELETE"); err != nil {
		t.Fatal(err)
	}
	if err := s.Authorize(2, "/api/v1/anything", "DELETE"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func testEnforcer(t *testing.T) *casbin.Enforcer {
	t.Helper()
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		t.Fatal(err)
	}
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatal(err)
	}
	return enforcer
}
