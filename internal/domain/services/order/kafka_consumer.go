package order

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	serviceErrors "wb-L0-task/internal/domain/errors"
	models "wb-L0-task/internal/domain/order"
	"wb-L0-task/internal/pkg/logger"
)

const phoneNumberRegex = `^(\+?\d{1,3})?\d{7,15}$`
const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

type KafkaConsumerService struct {
	storage Repository
}

func NewKafkaConsumerService(storage Repository) *KafkaConsumerService {
	return &KafkaConsumerService{
		storage: storage,
	}
}

func (s *KafkaConsumerService) SaveOrder(ctx context.Context, message []byte) error {
	var order *models.Order
	var err error
	if err = json.Unmarshal(message, &order); err != nil {
		logger.Error("Failed to unmarshal order", "error", err)
		return fmt.Errorf("invalid order format")
	}

	err = s.isValidOrder(order)
	if err != nil {
		return err
	}

	err = s.storage.Save(ctx, order)
	if err != nil {
		logger.Error("Failed to save order", "error", err)
		return fmt.Errorf("save order failed")
	}
	return nil
}

func (s *KafkaConsumerService) isValidOrder(order *models.Order) error {
	if order.UID == "" {
		return serviceErrors.ErrInvalidEntity.ForEntity("order_uid")
	}

	// Check correct phone number
	matched, err := regexp.MatchString(phoneNumberRegex, order.Delivery.Phone)
	if err != nil || !matched {
		return serviceErrors.ErrInvalidEntity.ForEntity("order.delivery.phone")
	}

	// Check correct email
	matched, err = regexp.MatchString(emailRegex, order.Delivery.Email)
	if err != nil || !matched {
		return serviceErrors.ErrInvalidEntity.ForEntity("order.delivery.email")
	}
	// Check total sum of items with payment
	goodsTotal := uint(0)
	for _, item := range order.Items {
		expectedTotal := item.Price * (100 - item.Sale) / 100
		if expectedTotal == item.TotalPrice {
			goodsTotal += item.TotalPrice
		} else {
			return serviceErrors.ErrInvalidEntity.ForEntity("order.item.total_price")
		}
	}
	if order.Payment.GoodsTotal != goodsTotal {
		return serviceErrors.ErrInvalidEntity.ForEntity("order.payment.goods_total")
	}
	if order.Payment.DeliveryCost+goodsTotal != order.Payment.Amount {
		return serviceErrors.ErrInvalidEntity.ForEntity("order.payment.goods_total")
	}

	return nil
}
