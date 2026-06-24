package iris

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/sylunbelievable/secondadmin/server/internal/config"
	"github.com/sylunbelievable/secondadmin/server/internal/service"
	"go.uber.org/zap"
)

type readinessCheck func(context.Context) error

const requestIDKey = "request_id"

func New(log *zap.Logger, ready readinessCheck, services *service.Services, cfg config.Config) *iris.Application {
	app := iris.New()
	app.UseRouter(recover.New())
	app.UseRouter(requestLog(log))
	app.UseRouter(cors(cfg.Auth.CORSOrigins))

	app.Get("/healthz", func(ctx iris.Context) { _ = ctx.JSON(iris.Map{"status": "ok"}) })
	app.Get("/readyz", func(ctx iris.Context) {
		checkCtx, cancel := context.WithTimeout(ctx.Request().Context(), 2*time.Second)
		defer cancel()
		if err := ready(checkCtx); err != nil {
			ctx.StatusCode(iris.StatusServiceUnavailable)
			_ = ctx.JSON(iris.Map{"status": "not_ready"})
			return
		}
		_ = ctx.JSON(iris.Map{"status": "ready"})
	})
	if cfg.Environment == "dev" {
		app.Get("/openapi.json", func(ctx iris.Context) {
			ctx.ContentType("application/json")
			_ = ctx.SendFile("docs/openapi.json", "")
		})
	}
	registerRoutes(app, services, cfg, log)
	return app
}

func requestLog(log *zap.Logger) iris.Handler {
	return func(ctx iris.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx.Header("X-Request-ID", requestID)
		ctx.Values().Set(requestIDKey, requestID)
		started := time.Now()
		ctx.Next()
		log.Info("http request",
			zap.String("request_id", requestID), zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()), zap.Int("status", ctx.GetStatusCode()),
			zap.Duration("duration", time.Since(started)),
		)
	}
}

func cors(origins []string) iris.Handler {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		allowed[origin] = true
	}
	return func(ctx iris.Context) {
		origin := ctx.GetHeader("Origin")
		if origin != "" && allowed[origin] {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Vary", "Origin")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token, X-Request-ID")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if ctx.Method() == "OPTIONS" {
			ctx.StatusCode(iris.StatusNoContent)
			return
		}
		ctx.Next()
	}
}
