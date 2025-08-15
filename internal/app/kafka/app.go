package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"log/slog"
	"wb-L0-task/internal/domain/services/order"
)

type App struct {
	logger   *slog.Logger
	consumer *kafka.Reader
	service  *order.KafkaConsumerService
}

func New(logger *slog.Logger, consumer *kafka.Reader, service *order.KafkaConsumerService) *App {
	return &App{
		logger:   logger,
		consumer: consumer,
		service:  service,
	}
}

func (a *App) Run(ctx context.Context) {
	a.logger.Info("Starting Kafka consumer...")
	for {
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}
		msg, err := a.consumer.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			a.logger.Error("Error while reading message", "err", err)
			continue
		}
		a.logger.Info("Message received",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"Key", string(msg.Key),
			"Value", string(msg.Value),
		)
		if err = a.consumer.CommitMessages(ctx, msg); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
		a.logger.Info("Received message", "kafka msg", string(msg.Value))
		err = a.service.SaveOrder(ctx, msg.Value)
		if err != nil {
			a.logger.Error("Failed to save order", "err", err)
		}
		continue
	}
}

func (a *App) Shutdown() {
	a.logger.Info("Shutting down Kafka consumer")
	if err := a.consumer.Close(); err != nil {
		a.logger.Error("Failed to close consumer", "err", err)
	}
}
