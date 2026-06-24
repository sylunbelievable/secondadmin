package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/kataras/iris/v12"
	"github.com/redis/go-redis/v9"
	"github.com/sylunbelievable/secondadmin/server/internal/config"
	"github.com/sylunbelievable/secondadmin/server/internal/repository"
	"github.com/sylunbelievable/secondadmin/server/internal/service"
	irisTransport "github.com/sylunbelievable/secondadmin/server/internal/transport/iris"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config   config.Config
	Log      *zap.Logger
	Database *gorm.DB
	sqlDB    *DatabaseHandle
	Redis    *redis.Client
	Casbin   *casbin.Enforcer
	watcher  interface{ Close() }
	HTTP     *iris.Application
	ready    atomic.Bool
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	log, err := newLogger(cfg.Log)
	if err != nil {
		return nil, err
	}

	db, sqlDB, err := OpenDatabase(ctx, cfg.Database)
	if err != nil {
		_ = log.Sync()
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		_ = sqlDB.Close()
		_ = log.Sync()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	policyDB := WriterDatabase(db)
	gormadapter.TurnOffAutoMigrate(policyDB)
	adapter, err := gormadapter.NewAdapterByDB(policyDB)
	if err != nil {
		_ = redisClient.Close()
		_ = sqlDB.Close()
		_ = log.Sync()
		return nil, fmt.Errorf("casbin adapter: %w", err)
	}
	enforcer, err := casbin.NewEnforcer("configs/casbin-model.conf", adapter)
	if err != nil {
		_ = redisClient.Close()
		_ = sqlDB.Close()
		_ = log.Sync()
		return nil, fmt.Errorf("casbin enforcer: %w", err)
	}
	watcher, err := rediswatcher.NewWatcher(cfg.Redis.Addr, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Database,
		},
		IgnoreSelf: true,
		OptionalUpdateCallback: func(string) {
			if err := enforcer.LoadPolicy(); err != nil {
				log.Error("reload casbin policy", zap.Error(err))
			}
		},
	})
	if err != nil {
		_ = redisClient.Close()
		_ = sqlDB.Close()
		_ = log.Sync()
		return nil, fmt.Errorf("casbin watcher: %w", err)
	}
	if err := enforcer.SetWatcher(watcher); err != nil {
		if closer, ok := watcher.(interface{ Close() }); ok {
			closer.Close()
		}
		_ = redisClient.Close()
		_ = sqlDB.Close()
		_ = log.Sync()
		return nil, fmt.Errorf("attach casbin watcher: %w", err)
	}

	app := &App{
		Config:   cfg,
		Log:      log,
		Database: db,
		sqlDB:    sqlDB,
		Redis:    redisClient,
		Casbin:   enforcer,
	}
	if closer, ok := watcher.(interface{ Close() }); ok {
		app.watcher = closer
	}
	services := service.New(repository.New(db), redisClient, enforcer, cfg)
	app.HTTP = irisTransport.New(log, app.Ready, services, cfg)
	app.ready.Store(true)
	return app, nil
}

func (a *App) Ready(ctx context.Context) error {
	if !a.ready.Load() {
		return errors.New("shutting down")
	}
	if err := a.sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	if err := a.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	return nil
}

func (a *App) Run() error {
	a.Log.Info("server starting", zap.String("addr", a.Config.HTTP.Addr), zap.String("env", a.Config.Environment))
	err := a.HTTP.Listen(a.Config.HTTP.Addr, iris.WithoutServerError(iris.ErrServerClosed))
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (a *App) Shutdown() error {
	a.ready.Store(false)
	ctx, cancel := context.WithTimeout(context.Background(), a.Config.HTTP.ShutdownTimeout)
	defer cancel()

	var result error
	if err := a.HTTP.Shutdown(ctx); err != nil {
		result = errors.Join(result, fmt.Errorf("http shutdown: %w", err))
	}
	if a.watcher != nil {
		a.watcher.Close()
	}
	if err := a.Redis.Close(); err != nil {
		result = errors.Join(result, fmt.Errorf("redis close: %w", err))
	}
	if err := a.sqlDB.Close(); err != nil {
		result = errors.Join(result, fmt.Errorf("database close: %w", err))
	}
	a.Log.Info("server stopped", zap.Duration("timeout", a.Config.HTTP.ShutdownTimeout))
	if err := a.Log.Sync(); err != nil {
		// stdout commonly returns EINVAL on sync; it is harmless at process exit.
		_ = err
	}
	return result
}
