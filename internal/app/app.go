package app

import (
	"context"
	"time"
	"wb-L0-task/internal/app/http"
	"wb-L0-task/internal/app/kafka"
	order_controller "wb-L0-task/internal/controllers/order"
	"wb-L0-task/internal/domain/order"
	order_service "wb-L0-task/internal/domain/services/order"
	"wb-L0-task/internal/pkg/cache"
	"wb-L0-task/internal/pkg/config"
	kafka_pkg "wb-L0-task/internal/pkg/kafka"
	"wb-L0-task/internal/pkg/logger"
	"wb-L0-task/internal/pkg/postgres"
	"wb-L0-task/internal/pkg/shutdown"
	repo_pkg "wb-L0-task/internal/repositories/postgres"
)

type App struct {
	HTTPApp  *http.App
	KafkaApp *kafka.App
}

func New(
	ctx context.Context,
	c *config.AppConfig,
) *App {
	pool, trManager, ctxGetter, err := postgres.SetupPostgres(ctx, c.Postgres)
	if err != nil {
		logger.Error("failed to setup postgres", "err", err)
	}

	orderRepo := repo_pkg.NewOrder(pool, trManager, ctxGetter)

	ordersCache := cache.NewCache[order.Order](c.Cache)
	orderService := order_service.New(ordersCache, orderRepo)
	err = orderService.InitCache(ctx)
	if err != nil {
		logger.Error("Failed to init cache", "err", err)
	}

	orderController := order_controller.New(orderService)

	consumer := kafka_pkg.NewConsumer(c.Kafka)

	httpApp := http.New(c, orderController)

	kafkaConsumerService := order_service.NewKafkaConsumerService(orderRepo)
	kafkaApp := kafka.New(consumer, kafkaConsumerService)

	shutdown.RegisterFn(func() {
		logger.Info("Shutting down")
		pool.Close()
		httpApp.Shutdown(time.Duration(c.Server.ShutdownTimeout))
		kafkaApp.Shutdown()
		logger.Info("Shutdown completed")
	})

	return &App{
		HTTPApp:  httpApp,
		KafkaApp: kafkaApp,
	}
}
