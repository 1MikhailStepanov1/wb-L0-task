package order

import (
	"context"
	"log/slog"
	serviceErrors "wb-L0-task/internal/domain/errors"
	model "wb-L0-task/internal/domain/order"
	"wb-L0-task/internal/repositories/postgres"
)

type Order struct {
	logger  *slog.Logger
	storage *postgres.Order
}

func New(logger *slog.Logger, storage *postgres.Order) *Order {
	return &Order{
		logger:  logger,
		storage: storage,
	}
}

func (o *Order) GetOrderById(ctx context.Context, orderId string) (*model.Order, error) {
	if exists, err := o.storage.Exists(ctx, orderId); exists && err == nil {
		res, err := o.storage.GetById(ctx, orderId)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else {
		return nil, serviceErrors.ErrNotFound.ForEntity("order")
	}
}
