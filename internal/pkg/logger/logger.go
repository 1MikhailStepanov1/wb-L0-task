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

var globalLogger = newDefault() //nolint: gochecknoglobals

func newDefault() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func New(cfg *Config) *slog.Logger {
	var log *slog.Logger
	switch cfg.LogMod {
	case devLogsMod:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case productionLogsMod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	globalLogger = log

	return log
}
