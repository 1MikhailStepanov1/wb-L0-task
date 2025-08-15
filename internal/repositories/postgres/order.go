package postgres

import (
	"context"
	"errors"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	models "wb-L0-task/internal/domain/order"
)

type Order struct {
	*Repo
}

func NewOrder(db *pgxpool.Pool, trManager TrManager, c *trmpgx.CtxGetter) *Order {
	return &Order{
		Repo: NewRepo(db, trManager, c),
	}
}

func (o *Order) GetById(ctx context.Context, orderUID string) (*models.Order, error) {
	var result models.Order
	err := o.trManager.Do(ctx, func(ctx context.Context) error {
		tx := o.getter.DefaultTrOrDB(ctx, o.db)
		err := tx.QueryRow(ctx, "SELECT * FROM orders WHERE uid = $1", orderUID).Scan(
			&result.UID,
			&result.TrackNumber,
			&result.Entry,
			&result.Locale,
			&result.InternalSignature,
			&result.CustomerID,
			&result.DeliveryService,
			&result.ShardKey,
			&result.StockManagementId,
			&result.OutOfFailureShard,
			&result.DateCreated,
		)
		if err != nil {
			return err
		}

		err = tx.QueryRow(ctx, "SELECT * FROM deliveries WHERE order_uid = $1", orderUID).Scan(
			&result.Delivery.ID,
			&result.Delivery.OrderUID,
			&result.Delivery.Name,
			&result.Delivery.Phone,
			&result.Delivery.Zip,
			&result.Delivery.City,
			&result.Delivery.Address,
			&result.Delivery.Region,
			&result.Delivery.Email,
		)
		if err != nil {
			return err
		}

		err = tx.QueryRow(ctx, "SELECT * FROM payments WHERE order_uid = $1", orderUID).Scan(
			&result.Payment.ID,
			&result.Payment.OrderUID,
			&result.Payment.TransactionID,
			&result.Payment.RequestID,
			&result.Payment.Currency,
			&result.Payment.Provider,
			&result.Payment.Amount,
			&result.Payment.PaymentDT,
			&result.Payment.Bank,
			&result.Payment.DeliveryCost,
			&result.Payment.GoodsTotal,
			&result.Payment.CustomFee,
		)
		if err != nil {
			return err
		}

		rows, err := tx.Query(ctx, "SELECT * FROM order_items WHERE order_uid = $1", orderUID)
		if err != nil {
			return fmt.Errorf("failed to get order items: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var item models.Item
			err = rows.Scan(
				&item.ID,
				&item.OrderUID,
				&item.ChartID,
				&item.TrackNumber,
				&item.Price,
				&item.RID,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NomenclatureID,
				&item.Brand,
				&item.Status,
			)
			if err != nil {
				return fmt.Errorf("failed to scan order item: %w", err)
			}
			result.Items = append(result.Items, item)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (o *Order) Exists(ctx context.Context, orderUID string) (bool, error) {
	var exists int8
	err := o.trManager.Do(ctx, func(ctx context.Context) error {
		tx := o.getter.DefaultTrOrDB(ctx, o.db)
		return tx.QueryRow(ctx, "SELECT 1 FROM orders WHERE uid = $1", orderUID).Scan(&exists)
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if order exists: %w", err)
	}
	return true, nil
}

func (o *Order) Save(ctx context.Context, order *models.Order) error {
	err := o.trManager.Do(ctx, func(ctx context.Context) error {
		tx := o.getter.DefaultTrOrDB(ctx, o.db)
		_, err := tx.Exec(ctx,
			`INSERT INTO orders(uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, 
                   shardkey, sm_id, oof_shard, date_created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			order.UID,
			order.TrackNumber,
			order.Entry,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.ShardKey,
			order.StockManagementId,
			order.OutOfFailureShard,
			order.DateCreated,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order: %w", err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO deliveries(id, order_uid, name, phone, zip, city, address, region, email)
				VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8)`,
			order.UID,
			order.Delivery.Name,
			order.Delivery.Phone,
			order.Delivery.Zip,
			order.Delivery.City,
			order.Delivery.Address,
			order.Delivery.Region,
			order.Delivery.Email,
		)
		if err != nil {
			return fmt.Errorf("failed to insert delivery order: %w", err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO payments(id, order_uid, transaction, request_id, currency, provider, amount, payment_dt, 
                     bank, delivery_cost, goods_total, custom_fee) 
					VALUES(gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			&order.UID,
			&order.Payment.TransactionID,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDT,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			return fmt.Errorf("failed to insert payment: %w", err)
		}

		batch := &pgx.Batch{}
		for _, item := range order.Items {
			batch.Queue(
				`INSERT INTO order_items(order_uid, chrt_id, track_number, price, rid, name, sale, 
                        size, total_price, nm_id, brand, status) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
				order.UID,
				&item.ChartID,
				&item.TrackNumber,
				&item.Price,
				&item.RID,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NomenclatureID,
				&item.Brand,
				&item.Status,
			)
		}
		batchRes := tx.SendBatch(ctx, batch)

		for range order.Items {
			_, err = batchRes.Exec()
			if err != nil {
				return fmt.Errorf("failed to insert order_items: %w", err)
			}
		}

		err = batchRes.Close()
		if err != nil {
			return fmt.Errorf("failed to close batch: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) GetOrders(ctx context.Context, limit int32) ([]models.Order, error) {
	orders := make([]models.Order, 0, limit)
	err := o.trManager.Do(ctx, func(ctx context.Context) error {
		tx := o.getter.DefaultTrOrDB(ctx, o.db)

		rows, err := tx.Query(ctx, "SELECT * FROM orders ORDER BY orders.date_created LIMIT $1", limit)
		defer rows.Close() //nolint:staticcheck
		if err != nil {
			return fmt.Errorf("failed to get orders: %w", err)
		}

		for rows.Next() {
			var order models.Order
			err = rows.Scan(
				&order.UID,
				&order.TrackNumber,
				&order.Entry,
				&order.Locale,
				&order.InternalSignature,
				&order.CustomerID,
				&order.DeliveryService,
				&order.ShardKey,
				&order.StockManagementId,
				&order.OutOfFailureShard,
				&order.DateCreated,
			)
			if err != nil {
				return fmt.Errorf("failed to scan order: %w", err)
			}
			orders = append(orders, order)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *Order) GetOrderDelivery(ctx context.Context, orderUID string) (*models.Delivery, error) {
	var delivery models.Delivery
	err := o.trManager.Do(ctx, func(ctx context.Context) error {
		tx := o.getter.DefaultTrOrDB(ctx, o.db)
		err := tx.QueryRow(ctx, "SELECT * FROM deliveries WHERE order_uid = $1", orderUID).Scan(
			&delivery.ID,
			&delivery.OrderUID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
		if err != nil {
			return fmt.Errorf("failed to get delivery: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (o *Order) GetOrderPayment(ctx context.Context, orderUID string) (*models.Payment, error) {
	var payment models.Payment
	err := o.trManager.Do(ctx, func(ctx context.Context) error {
		tx := o.getter.DefaultTrOrDB(ctx, o.db)
		err := tx.QueryRow(ctx, "SELECT * FROM payments WHERE order_uid = $1", orderUID).Scan(
			&payment.ID,
			&payment.OrderUID,
			&payment.TransactionID,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			return fmt.Errorf("failed to load order payment: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (o *Order) GetOrderItems(ctx context.Context, orderUID string) ([]models.Item, error) {
	var items []models.Item
	err := o.trManager.Do(ctx, func(ctx context.Context) error {
		tx := o.getter.DefaultTrOrDB(ctx, o.db)
		rows, err := tx.Query(ctx, "SELECT * FROM order_items WHERE order_uid = $1", orderUID)
		defer rows.Close() //nolint:staticcheck
		if err != nil {
			return fmt.Errorf("failed to load order items: %w", err)
		}

		for rows.Next() {
			var item models.Item
			err = rows.Scan(
				&item.ID,
				&item.OrderUID,
				&item.ChartID,
				&item.TrackNumber,
				&item.Price,
				&item.RID,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NomenclatureID,
				&item.Brand,
				&item.Status,
			)
			if err != nil {
				return fmt.Errorf("failed to load order item: %w", err)
			}
			items = append(items, item)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}
