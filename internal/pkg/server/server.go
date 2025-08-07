package server

import (
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	HTTPPort        int           `mapstructure:"http_port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

func New(c *Config) *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%d", c.HTTPPort),
	}
}
