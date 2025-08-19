package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"
	"time"
	"wb-L0-task/internal/controllers/order"
	"wb-L0-task/internal/pkg/config"
	"wb-L0-task/internal/pkg/logger"
	"wb-L0-task/internal/pkg/server"
)

type App struct {
	config     *config.AppConfig
	server     *http.Server
	controller *order.Controller
}

func New(
	config *config.AppConfig,
	controller *order.Controller,
) *App {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web"))))
	registerRoutes(r, controller)

	s := server.New(config.Server)
	s.Handler = r
	return &App{
		config:     config,
		server:     s,
		controller: controller,
	}
}

func (a *App) Run() {
	err := a.Start()
	if err != nil {
		logger.Error("Failed to start server", "err", err)
		os.Exit(1)
	}
}

func (a *App) Start() error {
	logger.Info("HTTP Server started")
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server error: %v", err)
	}
	return nil
}

func (a *App) Shutdown(shutdownTimeout time.Duration) {
	logger.Info("Stopping HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Error("HTTP shutdown with error", "err", err)
	}
}

func registerRoutes(router *chi.Mux, controller *order.Controller) {
	router.Get("/order/{order_uid}", controller.GetOrderById())
}
