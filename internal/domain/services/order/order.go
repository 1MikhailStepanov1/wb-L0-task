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
	// Cache search
	if order, exists := o.cache.Get(orderId); exists {
		return order.(*model.Order), nil
	}

	// Get it from DB
	if exists, err := o.storage.Exists(ctx, orderId); exists && err == nil { //TODO Update error handling from DB
		res, err := o.storage.GetById(ctx, orderId)
		if err != nil {
			return nil, err
		}
		// Saving order in cache
		o.cache.Set(orderId, res, 0)
		return res, nil
	} else {
		return nil, serviceErrors.ErrNotFound.ForEntity("order")
	}
}

func (o *Order) InitCache(ctx context.Context) error {
	o.logger.Info("Initializing orders cache", "cache_size", initCacheSize)
	orders, err := o.storage.GetOrders(ctx, initCacheSize)
	for _, order := range orders {
		delivery, err := o.storage.GetOrderDelivery(ctx, order.UID)
		if err != nil {
			return err
		}
		order.Delivery = *delivery

		payment, err := o.storage.GetOrderPayment(ctx, order.UID)
		if err != nil {
			return err
		}
		order.Payment = *payment

		items, err := o.storage.GetOrderItems(ctx, order.UID)
		if err != nil {
			return err
		}
		order.Items = items
	}
	if err != nil {
		return err
	}
	for _, order := range orders {
		o.cache.Set(order.UID, order, 5*time.Minute) //TODO Remove magic number
	}
	return nil
}
