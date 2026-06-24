package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string   `mapstructure:"-"`
	HTTP        HTTP     `mapstructure:"http"`
	Database    Database `mapstructure:"database"`
	Redis       Redis    `mapstructure:"redis"`
	Log         Log      `mapstructure:"log"`
	ID          ID       `mapstructure:"id"`
	Auth        Auth     `mapstructure:"auth"`
}

type HTTP struct {
	Addr            string        `mapstructure:"addr"`
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

type Database struct {
	Driver                string        `mapstructure:"driver"`
	DSN                   string        `mapstructure:"dsn"`
	ReadWrite             ReadWrite     `mapstructure:"readWrite"`
	MaxOpenConnections    int           `mapstructure:"maxOpenConnections"`
	MaxIdleConnections    int           `mapstructure:"maxIdleConnections"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connectionMaxLifetime"`
}

type ReadWrite struct {
	Enabled    bool     `mapstructure:"enabled"`
	WriterDSN  string   `mapstructure:"writerDsn"`
	ReaderDSNs []string `mapstructure:"readerDsns"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type ID struct {
	WorkerID int `mapstructure:"workerId"`
}

type Auth struct {
	JWTSecret       string        `mapstructure:"jwtSecret"`
	AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
	RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
	MaxDevices      int           `mapstructure:"maxDevices"`
	Cookie          Cookie        `mapstructure:"cookie"`
	CORSOrigins     []string      `mapstructure:"corsOrigins"`
}

type Cookie struct {
	Secure   bool   `mapstructure:"secure"`
	SameSite string `mapstructure:"sameSite"`
}

func Load() (Config, error) {
	env := envOr("APP_ENV", "dev")
	if !slices.Contains([]string{"dev", "staging", "prod"}, env) {
		return Config{}, fmt.Errorf("unsupported APP_ENV %q", env)
	}

	v := viper.New()
	v.SetConfigName("config." + env)
	v.SetConfigType("yaml")
	v.AddConfigPath(envOr("CONFIG_DIR", "configs"))
	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	bindings := map[string]string{
		"http.addr":                    "HTTP_ADDR",
		"database.driver":              "DB_DRIVER",
		"database.dsn":                 "DB_DSN",
		"database.readWrite.enabled":   "DB_READ_WRITE_ENABLED",
		"database.readWrite.writerDsn": "DB_WRITER_DSN",
		"redis.addr":                   "REDIS_ADDR",
		"redis.password":               "REDIS_PASSWORD",
		"id.workerId":                  "ID_WORKER_ID",
		"auth.jwtSecret":               "JWT_SECRET",
		"auth.maxDevices":              "AUTH_MAX_DEVICES",
		"auth.cookie.secure":           "COOKIE_SECURE",
	}
	for key, environmentVariable := range bindings {
		if err := v.BindEnv(key, environmentVariable); err != nil {
			return Config{}, fmt.Errorf("bind %s: %w", environmentVariable, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}
	if value := os.Getenv("DB_READER_DSNS"); value != "" {
		cfg.Database.ReadWrite.ReaderDSNs = strings.Split(value, ",")
	}
	cfg.Environment = env
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	switch {
	case c.HTTP.Addr == "":
		return errors.New("HTTP_ADDR is required")
	case c.Database.Driver != "postgres" && c.Database.Driver != "mysql":
		return errors.New("DB_DRIVER must be postgres or mysql")
	case c.Database.DSN == "":
		return errors.New("DB_DSN is required")
	case c.Database.ReadWrite.Enabled && len(cleanDSNs(c.Database.ReadWrite.ReaderDSNs)) == 0:
		return errors.New("database.readWrite.readerDsns is required when read-write splitting is enabled")
	case c.Redis.Addr == "":
		return errors.New("REDIS_ADDR is required")
	case c.Auth.JWTSecret == "":
		return errors.New("JWT_SECRET is required")
	case c.HTTP.ShutdownTimeout <= 0:
		return errors.New("http.shutdownTimeout must be positive")
	case c.Auth.AccessTokenTTL <= 0 || c.Auth.RefreshTokenTTL <= 0:
		return errors.New("token TTLs must be positive")
	case c.Auth.MaxDevices < 0:
		return errors.New("auth.maxDevices cannot be negative")
	case c.ID.WorkerID < 0 || c.ID.WorkerID > 1023:
		return errors.New("id.workerId must be between 0 and 1023")
	}

	if c.Environment == "staging" || c.Environment == "prod" {
		secret := strings.ToLower(strings.TrimSpace(c.Auth.JWTSecret))
		if len(c.Auth.JWTSecret) < 32 || slices.Contains([]string{"secret", "change_me", "changeme"}, secret) {
			return errors.New("JWT_SECRET must be at least 32 bytes and not a default value")
		}
		if !c.Auth.Cookie.Secure {
			return errors.New("auth.cookie.secure must be true outside dev")
		}
	}
	return nil
}

func (d Database) WriterDSN() string {
	if strings.TrimSpace(d.ReadWrite.WriterDSN) != "" {
		return d.ReadWrite.WriterDSN
	}
	return d.DSN
}

func (d Database) ReaderDSNs() []string {
	return cleanDSNs(d.ReadWrite.ReaderDSNs)
}

func cleanDSNs(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return result
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
