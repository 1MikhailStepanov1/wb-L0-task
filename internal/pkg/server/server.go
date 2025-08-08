package server

import (
	"fmt"
	"net/http"
)

type Config struct {
	HTTPPort        int `mapstructure:"http_port"`
	ShutdownTimeout int `mapstructure:"shutdown_timeout"`
}

func New(c *Config) *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%d", c.HTTPPort),
	}
}
