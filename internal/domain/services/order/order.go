package order

import (
	"context"
	"log/slog"
	"time"
	serviceErrors "wb-L0-task/internal/domain/errors"
	model "wb-L0-task/internal/domain/order"
	"wb-L0-task/internal/pkg/cache"
	"wb-L0-task/internal/repositories/postgres"
)

const initCacheSize = 10

type Order struct {
	logger  *slog.Logger
	storage *postgres.Order
	cache   *cache.Cache
}

func New(logger *slog.Logger, cache *cache.Cache, storage *postgres.Order) *Order {
	return &Order{
		logger:  logger,
		storage: storage,
		cache:   cache,
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

func (o *Order) InitCache() error {
	orders, err := o.storage.GetOrders(context.Background(), initCacheSize)
	if err != nil {
		return err
	}
	for _, order := range orders {
		o.cache.Set(order.UID, order, 5*time.Minute) //TODO Remove magic number
	}
	return nil
}
