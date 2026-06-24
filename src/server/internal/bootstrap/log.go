package bootstrap

import (
	"fmt"
	"os"

	"github.com/sylunbelievable/secondadmin/server/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(cfg config.Log) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("log level: %w", err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level)), nil
}
