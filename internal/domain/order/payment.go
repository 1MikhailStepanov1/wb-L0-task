package order

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID            uuid.UUID `json:"id" db:"id"`
	OrderUID      string    `json:"order_uid" db:"order_uid"`
	TransactionID string    `json:"transaction" db:"transaction"`
	RequestID     string    `json:"request_id" db:"request_id"`
	Currency      string    `json:"currency" db:"currency"`
	Provider      string    `json:"provider" db:"provider"`
	Amount        uint      `json:"amount" db:"amount"`
	PaymentDT     time.Time `json:"payment_dt" db:"payment_dt"`
	Bank          string    `json:"bank" db:"bank"`
	DeliveryCost  uint      `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal    uint      `json:"goods_total" db:"goods_total"`
	CustomFee     uint      `json:"custom_fee" db:"custom_fee"`
}
