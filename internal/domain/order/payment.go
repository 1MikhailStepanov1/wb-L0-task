package order

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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

func (p *Payment) MarshalJSON() ([]byte, error) {
	type Alias Payment

	return json.Marshal(&struct {
		PaymentDT int64 `json:"payment_dt"`
		*Alias
	}{
		PaymentDT: p.PaymentDT.Unix(),
		Alias:     (*Alias)(p),
	})
}

type Payment struct {
	ID            uuid.UUID `json:"-"             db:"id"`
	OrderUID      string    `json:"-"             db:"order_uid"`
	TransactionID string    `json:"transaction"   db:"transaction"`
	RequestID     string    `json:"request_id"    db:"request_id"`
	Currency      string    `json:"currency"      db:"currency"`
	Provider      string    `json:"provider"      db:"provider"`
	Amount        uint      `json:"amount"        db:"amount"`
	PaymentDT     time.Time `json:"payment_dt"    db:"payment_dt"`
	Bank          string    `json:"bank"          db:"bank"`
	DeliveryCost  uint      `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal    uint      `json:"goods_total"   db:"goods_total"`
	CustomFee     uint      `json:"custom_fee"    db:"custom_fee"`
}
