package iris

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sylunbelievable/secondadmin/server/internal/entity"
	"github.com/sylunbelievable/secondadmin/server/internal/service"
	"go.uber.org/zap"
)

const principalKey = "principal"

func authenticate(services *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var token string
		if header := ctx.GetHeader("Authorization"); header != "" {
			var ok bool
			token, ok = strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				writeError(ctx, service.ErrUnauthorized)
				return
			}
		} else {
			token = ctx.GetCookie("access_token")
		}
		principal, err := services.Authenticate(ctx.Request().Context(), token)
		if err != nil {
			writeError(ctx, err)
			return
		}
		if principal.AuthMode == "cookie" && isWrite(ctx.Method()) {
			if err := services.CheckCSRF(ctx.Request().Context(), principal.SessionID, ctx.GetHeader("X-CSRF-Token")); err != nil {
				writeError(ctx, err)
				return
			}
		}
		ctx.Values().Set(principalKey, principal)
		ctx.Next()
	}
}

func authorize(services *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		p := principal(ctx)
		if err := services.Authorize(p.UserID, ctx.Path(), ctx.Method()); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.Next()
	}
}

func audit(services *service.Services, log *zap.Logger) iris.Handler {
	return func(ctx iris.Context) {
		if !isWrite(ctx.Method()) {
			ctx.Next()
			return
		}
		started := time.Now()
		ctx.Next()
		p := principal(ctx)
		if p.UserID == 0 {
			return
		}
		requestID, _ := ctx.Values().Get(requestIDKey).(string)
		item := entity.OperationLog{
			UserID: p.UserID, RequestID: requestID, Method: ctx.Method(), Path: ctx.Path(),
			StatusCode: ctx.GetStatusCode(), DurationMS: time.Since(started).Milliseconds(),
			IP: ctx.RemoteAddr(), UserAgent: ctx.GetHeader("User-Agent"),
		}
		if err := services.AddOperationLog(context.Background(), &item); err != nil {
			log.Error("write operation log", zap.Error(err))
		}
	}
}

func principal(ctx iris.Context) service.Principal {
	value := ctx.Values().Get(principalKey)
	p, _ := value.(service.Principal)
	return p
}

func isWrite(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE"
}

func writeError(ctx iris.Context, err error) {
	status, code := iris.StatusInternalServerError, "INTERNAL_ERROR"
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		status, code = iris.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS"
	case errors.Is(err, service.ErrUnauthorized):
		status, code = iris.StatusUnauthorized, "AUTH_UNAUTHORIZED"
	case errors.Is(err, service.ErrForbidden):
		status, code = iris.StatusForbidden, "AUTH_FORBIDDEN"
	case errors.Is(err, service.ErrInvalidInput):
		status, code = iris.StatusBadRequest, "INVALID_INPUT"
	case errors.Is(err, service.ErrNotFound):
		status, code = iris.StatusNotFound, "NOT_FOUND"
	case errors.Is(err, service.ErrConflict):
		status, code = iris.StatusConflict, "CONFLICT"
	case errors.Is(err, service.ErrDependency):
		status, code = iris.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE"
	case errors.Is(err, service.ErrRateLimited):
		status, code = iris.StatusTooManyRequests, "AUTH_RATE_LIMITED"
	}
	ctx.StatusCode(status)
	requestID, _ := ctx.Values().Get(requestIDKey).(string)
	_ = ctx.JSON(iris.Map{"code": code, "message": err.Error(), "requestId": requestID})
}
