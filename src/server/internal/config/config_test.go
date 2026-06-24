package config

import "testing"

func TestProductionRejectsWeakJWTSecret(t *testing.T) {
	cfg := Config{
		Environment: "prod",
		HTTP:        HTTP{Addr: ":8080", ShutdownTimeout: 10},
		Database:    Database{Driver: "postgres", DSN: "dsn"},
		Redis:       Redis{Addr: "redis:6379"},
		Auth: Auth{
			JWTSecret:       "change_me",
			AccessTokenTTL:  1,
			RefreshTokenTTL: 1,
			Cookie:          Cookie{Secure: true},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected weak production JWT secret to be rejected")
	}
}

func TestReadWriteSplittingDefaultsToOff(t *testing.T) {
	cfg := Config{
		Environment: "dev",
		HTTP:        HTTP{Addr: ":8080", ShutdownTimeout: 10},
		Database:    Database{Driver: "postgres", DSN: "writer"},
		Redis:       Redis{Addr: "redis:6379"},
		Auth: Auth{
			JWTSecret:       "dev",
			AccessTokenTTL:  1,
			RefreshTokenTTL: 1,
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if got := cfg.Database.WriterDSN(); got != "writer" {
		t.Fatalf("writer dsn = %q", got)
	}
}

func TestReadWriteSplittingRequiresReadersWhenEnabled(t *testing.T) {
	cfg := Config{
		Environment: "dev",
		HTTP:        HTTP{Addr: ":8080", ShutdownTimeout: 10},
		Database:    Database{Driver: "postgres", DSN: "writer", ReadWrite: ReadWrite{Enabled: true}},
		Redis:       Redis{Addr: "redis:6379"},
		Auth: Auth{
			JWTSecret:       "dev",
			AccessTokenTTL:  1,
			RefreshTokenTTL: 1,
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected enabled read-write splitting to require readers")
	}
}
