package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/sylunbelievable/secondadmin/server/internal/bootstrap"
	"github.com/sylunbelievable/secondadmin/server/internal/config"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	app, err := bootstrap.New(ctx, cfg)
	if err != nil {
		return err
	}

	runErr := make(chan error, 1)
	go func() { runErr <- app.Run() }()

	var serverErr error
	select {
	case err := <-runErr:
		if err != nil {
			app.Log.Error("server failed", zap.Error(err))
			serverErr = fmt.Errorf("run server: %w", err)
		}
	case <-ctx.Done():
	}
	return errors.Join(serverErr, app.Shutdown())
}
