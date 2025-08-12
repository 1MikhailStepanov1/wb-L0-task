package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
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
		select {
		case <-ctx.Done():
			a.logger.Info("Stopping Kafka consumer...")
			if err := a.consumer.Close(); err != nil {
				a.logger.Error("Failed to close consumer", "err", err)
				return
			}
			a.logger.Info("Kafka consumer stopped successfully")
			return
		default:
			a.logger.Info("Reading message...")
			//msg, err := a.consumer.ReadMessage(ctx)
			//if err != nil {
			//	a.logger.Error("Error while reading message", "err", err)
			//	continue
			//}
			msg, err := a.consumer.FetchMessage(ctx)
			if err != nil {
				a.logger.Error("Error while reading message", "err", err)
				continue
			}
			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			if err = a.consumer.CommitMessages(ctx, msg); err != nil {
				log.Fatal("failed to commit messages:", err)
			}
			a.logger.Info("Received message", string(msg.Value))
			err = a.service.SaveOrder(ctx, msg.Value)
			if err != nil {
				a.logger.Error("Failed to save order", "err", err)
			}
			continue
		}
	}
}

func (a *App) Shutdown() error {
	err := a.consumer.Close()
	if err != nil {
		a.logger.Error("Failed to close consumer", err)
		return err
	}
	return nil
}
