package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	models "wb-L0-task/internal/domain/order"
	repo_pkg "wb-L0-task/internal/repositories/postgres"
)

type Repository interface {
	SaveOrder(ctx context.Context, order *Order) error
	SaveDelivery(ctx context.Context, delivery models.Delivery) error
	SavePayment(ctx context.Context, payment *models.Payment) error
	SaveItems(ctx context.Context, items []models.Item) error
}

type KafkaConsumerService struct {
	logger  *slog.Logger
	storage *repo_pkg.Order
}

func NewKafkaConsumerService(logger *slog.Logger, storage *repo_pkg.Order) *KafkaConsumerService {
	return &KafkaConsumerService{
		logger:  logger,
		storage: storage,
	}
}

func (s *KafkaConsumerService) SaveOrder(ctx context.Context, message []byte) error {
	var order *models.Order
	var err error
	if err = json.Unmarshal(message, &order); err != nil {
		s.logger.Error("Failed to unmarshal order", "error", err)
		return fmt.Errorf("invalid order format")
	}

	if order.UID == "" {
		return fmt.Errorf("order UID is required")
	}

	err = s.storage.Save(ctx, order)
	if err != nil {
		s.logger.Error("Failed to save order", "error", err)
		return fmt.Errorf("save order failed")
	}
	return nil
}
