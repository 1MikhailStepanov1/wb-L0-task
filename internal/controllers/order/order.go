package order

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	serviceErrors "wb-L0-task/internal/domain/errors"
	model "wb-L0-task/internal/domain/order"
	"wb-L0-task/internal/pkg/logger"
)

type Service interface {
	GetOrderById(ctx context.Context, orderId string) (*model.Order, error)
}

type Controller struct {
	service Service
}

func New(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) GetOrderById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := chi.URLParam(r, "order_uid")
		if orderUID == "" {
			http.Error(w, "order_uid is required", http.StatusBadRequest)
			return
		}

		order, err := c.service.GetOrderById(r.Context(), orderUID)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(order); err != nil {
			logger.Error("Failed to encode response", "err", err)
		}
	}
}
