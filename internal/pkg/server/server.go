package server

import (
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	HTTPPort              int16 `mapstructure:"http_port"`
	ShutdownTimeout       int16 `mapstructure:"shutdown_timeout"`
	HTTPReadTimeout       int16 `mapstructure:"http_read_timeout"`
	HTTPWriteTimeout      int16 `mapstructure:"http_write_timeout"`
	HTTPIdleTimeout       int16 `mapstructure:"http_idle_timeout"`
	HTTPReadHeaderTimeout int16 `mapstructure:"http_read_header_timeout"`
}

func New(c *Config) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", c.HTTPPort),
		ReadTimeout:       time.Duration(c.HTTPReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(c.HTTPWriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(c.HTTPIdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(c.HTTPReadHeaderTimeout) * time.Second,
	}
}
