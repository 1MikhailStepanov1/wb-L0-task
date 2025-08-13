package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"wb-L0-task/internal/controllers/order"
	"wb-L0-task/internal/pkg/config"
	"wb-L0-task/internal/pkg/server"
)

type App struct {
	logger     *slog.Logger
	config     *config.AppConfig
	server     *http.Server
	controller *order.Controller
}

func New(
	logger *slog.Logger,
	config *config.AppConfig,
	controller *order.Controller,
) *App {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	registerRoutes(r, controller)

	s := server.New(config.Server)
	s.Handler = r
	return &App{
		logger:     logger,
		config:     config,
		server:     s,
		controller: controller,
	}
}

// TODO Error handling
func (a *App) Run() {
	err := a.Start()
	if err != nil {
		panic(err)
	}
}

func (a *App) Start() error {
	a.logger.Info("HTTP Server started")
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server error: %v", err)
	}
	return nil
}

// TODO Graceful shutdown
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("Stopping HTTP server")
	return a.server.Shutdown(ctx)
}

func registerRoutes(router *chi.Mux, controller *order.Controller) {
	router.Get("/order/{order_uid}", controller.GetOrderById())
}
