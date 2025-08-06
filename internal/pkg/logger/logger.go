package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	LogMod string `mapstructure:"mod"`
}

const (
	devLogsMod        = "DEV"
	productionLogsMod = "PROD"
)

func New(cfg Config) *slog.Logger {
	var log *slog.Logger
	switch cfg.LogMod {
	case devLogsMod:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case productionLogsMod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
