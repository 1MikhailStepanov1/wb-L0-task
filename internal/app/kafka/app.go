package kafka

import (
	"context"
	"errors"
	"io"
	"log"

	"wb-L0-task/internal/domain/services/order"
	"wb-L0-task/internal/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type App struct {
	consumer *kafka.Reader
	service  *order.KafkaConsumerService
}

func New(consumer *kafka.Reader, service *order.KafkaConsumerService) *App {
	return &App{
		consumer: consumer,
		service:  service,
	}
}

func (a *App) Run(ctx context.Context) {
	logger.Info("Starting Kafka consumer...")
	for {
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}
		msg, err := a.consumer.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("Error while reading message", "err", err)
			continue
		}
		logger.Info("Message received",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"Key", string(msg.Key),
			"Value", string(msg.Value),
		)
		if err = a.consumer.CommitMessages(ctx, msg); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
		logger.Info("Received message", "kafka msg", string(msg.Value))
		err = a.service.SaveOrder(ctx, msg.Value)
		if err != nil {
			logger.Error("Failed to save order", "err", err)
		}
		continue
	}
}

func (a *App) Shutdown() {
	logger.Info("Shutting down Kafka consumer")
	if err := a.consumer.Close(); err != nil {
		logger.Error("Failed to close consumer", "err", err)
	}
}
