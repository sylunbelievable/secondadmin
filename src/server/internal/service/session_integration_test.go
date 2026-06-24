package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sylunbelievable/secondadmin/server/internal/config"
)

func TestSessionLimitAndRotation(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR is not set")
	}
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: addr})
	t.Cleanup(func() {
		keys, _ := client.Keys(ctx, "auth:*").Result()
		if len(keys) > 0 {
			_ = client.Del(ctx, keys...).Err()
		}
		_ = client.Close()
	})
	s := Services{Redis: client, Config: config.Config{Auth: config.Auth{RefreshTokenTTL: time.Minute, MaxDevices: 1}}}
	if err := s.createSession(ctx, 1, "old", "a", "bearer", "old.token", "csrf"); err != nil {
		t.Fatal(err)
	}
	if err := s.createSession(ctx, 1, "new", "b", "bearer", "new.token", "csrf"); err != nil {
		t.Fatal(err)
	}
	if client.Exists(ctx, sessionKey("old")).Val() != 0 {
		t.Fatal("oldest session was not evicted")
	}
	rotated, err := rotateSessionScript.Run(ctx, client, []string{sessionKey("new")},
		tokenHash("new.token"), tokenHash("next.token"), tokenHash("next-csrf"), 60).Int()
	if err != nil || rotated != 1 {
		t.Fatal("refresh token was not rotated")
	}
	replayed, _ := rotateSessionScript.Run(ctx, client, []string{sessionKey("new")},
		tokenHash("new.token"), tokenHash("again.token"), tokenHash("again-csrf"), 60).Int()
	if replayed != 0 {
		t.Fatal("old refresh token was accepted after rotation")
	}
}
