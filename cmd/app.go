package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"wb-L0-task/internal/app"
	"wb-L0-task/internal/pkg/config"
	"wb-L0-task/internal/pkg/logger"
	"wb-L0-task/internal/pkg/shutdown"
)

func makeQuitSignal() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func main() {
	cfg, err := config.New()
	if err != nil {

	}
	_ = logger.New(cfg.Logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.New(ctx, cfg)

	go func() {
		application.HTTPApp.Run()
	}()

	go func() {
		application.KafkaApp.Run(ctx)
	}()

	shutdown.WaitSignal(makeQuitSignal())
}
