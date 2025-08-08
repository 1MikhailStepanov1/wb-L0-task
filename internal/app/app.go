package app

import (
	"context"
	"log/slog"
	"wb-L0-task/internal/app/http"
	order_controller "wb-L0-task/internal/controllers/order"
	order_service "wb-L0-task/internal/domain/services/order"
	"wb-L0-task/internal/pkg/config"
	"wb-L0-task/internal/pkg/postgres"
	repo_pkg "wb-L0-task/internal/repositories/postgres"
)

type App struct {
	HTTPApp *http.App
}

// TODO Handling panic
func New(
	ctx context.Context,
	c *config.AppConfig,
	logger *slog.Logger,
) *App {
	pool, trManager, ctxGetter, err := postgres.SetupPostgres(ctx, c.Postgres)
	if err != nil {
		panic(err)
	}

	orderRepo := repo_pkg.NewOrder(pool, trManager, ctxGetter)
	orderService := order_service.New(logger, orderRepo)
	orderController := order_controller.New(logger, orderService)

	httpApp := http.New(logger, c, orderController)
	return &App{
		HTTPApp: httpApp,
	}
}
