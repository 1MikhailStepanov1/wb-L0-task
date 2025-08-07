package postgres

import (
	"context"
	"errors"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	model "wb-L0-task/internal/domain/order"
)

type Order struct {
	*Repo
}

func NewOrder(db *pgxpool.Pool, trManager TrManager, c *trmpgx.CtxGetter) *Order {
	return &Order{
		Repo: NewRepo(db, trManager, c),
	}
}

func (o *Order) GetById(ctx context.Context, orderUID string) (*model.Order, error) {
	var result model.Order
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
			var item model.Item
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
	var exists bool
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
	return exists, nil
}

func (o *Order) Save(ctx context.Context, order model.Order) error {
	return nil
}
