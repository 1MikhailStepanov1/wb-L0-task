package order

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	model "wb-L0-task/internal/domain/order"
)

type Service interface {
	GetOrderById(ctx context.Context, orderId string) (*model.Order, error)
}

type Controller struct {
	service Service
}

func New(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetOrderById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := chi.URLParam(r, "order_uid")
		if orderUID != "" {
			http.Error(w, "order_uid is required", http.StatusBadRequest)
			return

		}

		order, err := c.service.GetOrderById(r.Context(), orderUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(order)
	}
}
