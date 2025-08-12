package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb-L0-task/internal/app"
	"wb-L0-task/internal/pkg/config"
	"wb-L0-task/internal/pkg/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	l := logger.New(cfg.Logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.New(ctx, cfg, l)

	go func() {
		application.HTTPApp.Run()
	}()

	go func() {
		application.KafkaApp.Run(ctx)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	l.Info("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer shutdownCancel()

	if err = application.HTTPApp.Stop(shutdownCtx); err != nil {
		l.Error("HTTP shutdown with error", err)
	} else {
		l.Info("HTTP server shut down cleanly")
	}

	if err = application.KafkaApp.Shutdown(); err != nil {
		l.Error("Kafka shutdown with error", err)
	} else {
		l.Info("Kafka shutdown cleanly")
	}
}
