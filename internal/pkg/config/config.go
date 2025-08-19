package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"wb-L0-task/internal/pkg/cache"
	"wb-L0-task/internal/pkg/kafka"
	"wb-L0-task/internal/pkg/logger"
	"wb-L0-task/internal/pkg/postgres"
	"wb-L0-task/internal/pkg/server"
)

const defaultConfigFileName = "config.yaml"

var ErrEmptyPath = errors.New("path to config must not be empty")

type AppConfig struct {
	Cache    *cache.Config
	Kafka    *kafka.Config
	Logger   *logger.Config
	Postgres *postgres.Config
	Server   *server.Config
}

func New() (*AppConfig, error) {
	return NewFromFilePath(defaultConfigFileName)
}

func NewFromFilePath(path string) (*AppConfig, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}

	cfg, err := initConfig(path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func initConfig(configFile string) (*AppConfig, error) {
	viperInstance := viper.New()
	ext := strings.TrimLeft(filepath.Ext(configFile), ".")
	viperInstance.SetConfigFile(configFile)
	viperInstance.SetConfigType(ext)

	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("v.ReadInConfig: %w", err)
	}

	err = godotenv.Load()
	if err != nil {
		logger.Debug("Can`t load environment variables from file", "err", err)
	}

	for _, key := range viperInstance.AllKeys() {
		value := viperInstance.GetString(key)
		if value == "" {
			continue
		}

		viperInstance.Set(key, os.ExpandEnv(value))
	}

	cfg := new(AppConfig)
	if err = viperInstance.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("v.Unmarshal: %w", err)
	}

	return cfg, nil
}
