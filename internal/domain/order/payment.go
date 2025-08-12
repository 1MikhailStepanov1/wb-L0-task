package order

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

func (p *Payment) UnmarshalJSON(data []byte) error {
	type Alias Payment
	aux := &struct {
		PaymentDT int64 `json:"payment_dt"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.PaymentDT = time.Unix(aux.PaymentDT, 0)
	return nil
}

type Payment struct {
	ID            uuid.UUID `json:"-" db:"id"`
	OrderUID      string    `json:"-" db:"order_uid"`
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
