package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var symbolsMap = map[int]string{
	1: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+ ",
	2: "0123456789",
	3: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ",
}

type Order struct {
	UID               string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	StockManagementId int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OutOfFailureShard string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	TransactionID string `json:"transaction"`
	RequestID     string `json:"request_id"`
	Currency      string `json:"currency"`
	Provider      string `json:"provider"`
	Amount        uint   `json:"amount"`
	PaymentDT     int64  `json:"payment_dt"`
	Bank          string `json:"bank"`
	DeliveryCost  uint   `json:"delivery_cost"`
	GoodsTotal    uint   `json:"goods_total"`
	CustomFee     uint   `json:"custom_fee"`
}

type Item struct {
	ChartID        int64  `json:"chrt_id"`
	TrackNumber    string `json:"track_number"`
	Price          uint   `json:"price"`
	RID            string `json:"rid"`
	Name           string `json:"name"`
	Sale           uint   `json:"sale"`
	Size           string `json:"size"`
	TotalPrice     uint   `json:"total_price"`
	NomenclatureID int64  `json:"nm_id"`
	Brand          string `json:"brand"`
	Status         int    `json:"status"`
}

func main() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "orders", 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		order := generateOrder()
		orderJson, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.WriteMessages(kafka.Message{
			Key:   []byte(order.UID),
			Value: orderJson,
		})
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				log.Printf("Work deadline exceeded")
				return
			}
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func generateDelivery() *Delivery {
	return &Delivery{
		Name:    generateRandomString(1, rand.Intn(100)),
		Phone:   fmt.Sprintf("+%s", generateRandomString(2, rand.Intn(18))),
		Zip:     strconv.Itoa(rand.Intn(1000000)),
		City:    generateRandomString(3, rand.Intn(30)),
		Address: generateRandomString(3, rand.Intn(50)),
		Region:  generateRandomString(3, rand.Intn(25)),
		Email: fmt.Sprintf("%s@%s.%s",
			generateRandomString(3, rand.Intn(15)),
			generateRandomString(3, rand.Intn(15)),
			generateRandomString(3, rand.Intn(15)),
		),
	}
}

func generatePayment(goodsTotal uint) *Payment {
	deliveryCost := uint(rand.Intn(10000))
	customFee := uint(rand.Intn(1000))
	amount := goodsTotal + customFee + deliveryCost
	return &Payment{
		TransactionID: generateRandomString(3, rand.Intn(15)),
		RequestID:     generateRandomString(3, rand.Intn(15)),
		Currency:      generateRandomString(3, rand.Intn(3)),
		Provider:      generateRandomString(3, rand.Intn(10)),
		Amount:        amount,
		PaymentDT:     time.Now().Add(time.Duration(rand.Intn(1000)) * time.Minute).Unix(),
		Bank:          generateRandomString(3, rand.Intn(15)),
		DeliveryCost:  deliveryCost,
		GoodsTotal:    goodsTotal,
		CustomFee:     customFee,
	}
}

func generateOrderItem() *Item {
	price := uint(rand.Intn(10000000))
	sale := uint(rand.Intn(100))
	totalPrice := price * (100 - sale) / 100

	return &Item{
		ChartID:        rand.Int63n(10000000),
		TrackNumber:    generateRandomString(1, rand.Intn(15)),
		Price:          price,
		RID:            generateRandomString(3, rand.Intn(25)),
		Name:           generateRandomString(3, rand.Intn(50)),
		Sale:           sale,
		Size:           generateRandomString(2, rand.Intn(5)),
		TotalPrice:     totalPrice,
		NomenclatureID: rand.Int63n(10000000),
		Brand:          generateRandomString(1, rand.Intn(50)),
		Status:         rand.Intn(1000),
	}
}

func generateOrder() *Order {
	itemsLen := rand.Intn(25)
	items := make([]Item, itemsLen)
	var goodsTotal uint = 0
	for i := range items {
		items[i] = *generateOrderItem()
		goodsTotal += items[i].TotalPrice
	}
	return &Order{
		UID:               generateRandomString(3, 15),
		TrackNumber:       generateRandomString(1, rand.Intn(15)),
		Entry:             generateRandomString(3, rand.Intn(10)),
		Delivery:          *generateDelivery(),
		Payment:           *generatePayment(goodsTotal),
		Items:             items,
		Locale:            generateRandomString(3, rand.Intn(2)),
		InternalSignature: generateRandomString(1, rand.Intn(100)),
		CustomerID:        generateRandomString(3, rand.Intn(25)),
		DeliveryService:   generateRandomString(3, rand.Intn(15)),
		ShardKey:          generateRandomString(2, rand.Intn(5)),
		StockManagementId: rand.Intn(10000),
		DateCreated:       time.Now().Add(time.Duration(rand.Intn(1000)) * time.Minute),
		OutOfFailureShard: generateRandomString(2, rand.Intn(5)),
	}
}

func generateRandomString(symbolPack int, length int) string {
	res := make([]byte, length)
	symbols := symbolsMap[symbolPack]
	for i := range res {
		res[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(res)
}
