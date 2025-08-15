package app

import (
	"context"
	"log/slog"
	"wb-L0-task/internal/app/http"
	"wb-L0-task/internal/app/kafka"
	order_controller "wb-L0-task/internal/controllers/order"
	"wb-L0-task/internal/domain/order"
	order_service "wb-L0-task/internal/domain/services/order"
	"wb-L0-task/internal/pkg/cache"
	"wb-L0-task/internal/pkg/config"
	kafka_pkg "wb-L0-task/internal/pkg/kafka"
	"wb-L0-task/internal/pkg/postgres"
	repo_pkg "wb-L0-task/internal/repositories/postgres"
)

type App struct {
	HTTPApp  *http.App
	KafkaApp *kafka.App
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

	ordersCache := cache.NewCache[order.Order](c.Cache)
	orderService := order_service.New(logger, ordersCache, orderRepo)
	err = orderService.InitCache(ctx)
	if err != nil {
		logger.Error("Failed to init cache", "err", err)
	}

	orderController := order_controller.New(logger, orderService)

	consumer, err := kafka_pkg.NewConsumer(c.Kafka)
	if err != nil {
		panic(err)
	}

	httpApp := http.New(logger, c, orderController)

	kafkaConsumerService := order_service.NewKafkaConsumerService(logger, orderRepo)
	kafkaApp := kafka.New(logger, consumer, kafkaConsumerService)
	return &App{
		HTTPApp:  httpApp,
		KafkaApp: kafkaApp,
	}
}
