package order

import (
	"encoding/json"
	"time"
)

func (o *Order) UnmarshalJSON(data []byte) error {
	type Alias Order
	aux := &struct {
		DateCreated string `json:"date_created"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.DateCreated != "" {
		t, err := time.Parse(time.RFC3339, aux.DateCreated)
		if err != nil {
			return err
		}
		o.DateCreated = t
	}
	return nil
}

type Order struct {
	UID               string    `json:"order_uid" db:"uid"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale" db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	ShardKey          string    `json:"shardkey" db:"shardkey"`
	StockManagementId int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OutOfFailureShard string    `json:"oof_shard" db:"oof_shard"`
}

type Item struct {
	ID             int    `json:"-" db:"id"`
	OrderUID       string `json:"-" db:"order_uid"`
	ChartID        int64  `json:"chrt_id" db:"chrt_id"`
	TrackNumber    string `json:"track_number" db:"track_number"`
	Price          uint   `json:"price" db:"price"`
	RID            string `json:"rid" db:"rid"`
	Name           string `json:"name" db:"name"`
	Sale           int8   `json:"sale" db:"sale"`
	Size           string `json:"size" db:"size"`
	TotalPrice     uint   `json:"total_price" db:"total_price"`
	NomenclatureID int64  `json:"nm_id" db:"nm_id"`
	Brand          string `json:"brand" db:"brand"`
	Status         int    `json:"status" db:"status"`
}
