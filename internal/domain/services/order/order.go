package order

import (
	"context"
	"time"
	serviceErrors "wb-L0-task/internal/domain/errors"
	model "wb-L0-task/internal/domain/order"
	"wb-L0-task/internal/pkg/cache"
	"wb-L0-task/internal/pkg/logger"
)

const initCacheSize = 10

type Repository interface {
	GetById(ctx context.Context, orderUID string) (*model.Order, error)
	Exists(ctx context.Context, orderUID string) (bool, error)
	Save(ctx context.Context, order *model.Order) error
	GetOrders(ctx context.Context, limit int32) ([]model.Order, error)
	GetOrderDelivery(ctx context.Context, orderUID string) (*model.Delivery, error)
	GetOrderPayment(ctx context.Context, orderUID string) (*model.Payment, error)
	GetOrderItems(ctx context.Context, orderUID string) ([]model.Item, error)
}

type Order struct {
	storage Repository
	cache   *cache.Cache[model.Order]
}

func New(cache *cache.Cache[model.Order], storage Repository) *Order {
	return &Order{
		storage: storage,
		cache:   cache,
	}
}

func (o *Order) GetOrderById(ctx context.Context, orderId string) (*model.Order, error) {
	// Cache search
	if order, exists := o.cache.Get(orderId); exists {
		return &order, nil
	}

	// Get it from DB
	if exists, err := o.storage.Exists(ctx, orderId); exists && err == nil { //TODO Update error handling from DB
		res, err := o.storage.GetById(ctx, orderId)
		if err != nil {
			return nil, err
		}
		// Saving order in cache
		o.cache.Set(orderId, *res, 0)
		return res, nil
	} else {
		return nil, serviceErrors.ErrNotFound.ForEntity("order")
	}
}

func (o *Order) InitCache(ctx context.Context) error {
	logger.Info("Initializing orders cache", "cache_size", initCacheSize)
	orders, err := o.storage.GetOrders(ctx, initCacheSize)
	for i := range orders {
		delivery, err := o.storage.GetOrderDelivery(ctx, orders[i].UID)
		if err != nil {
			logger.Error("Failed to get order delivery", "order_uid", orders[i].UID, "error", err)
			return err
		}
		orders[i].Delivery = *delivery

		payment, err := o.storage.GetOrderPayment(ctx, orders[i].UID)
		if err != nil {
			logger.Error("Failed to get order payment", "order_uid", orders[i].UID, "error", err)
			return err
		}
		orders[i].Payment = *payment

		items, err := o.storage.GetOrderItems(ctx, orders[i].UID)
		if err != nil {
			logger.Error("Failed to get order items", "order_uid", orders[i].UID, "error", err)
			return err
		}
		orders[i].Items = items
	}
	if err != nil {
		return err
	}
	for _, order := range orders {
		o.cache.Set(order.UID, order, 5*time.Minute) //TODO Remove magic number
	}
	return nil
}
