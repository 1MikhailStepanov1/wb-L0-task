// internal/domain/order/controller_test.go
package order

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	serviceErrors "wb-L0-task/internal/domain/errors"
	model "wb-L0-task/internal/domain/order"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOrderById_Success(t *testing.T) {
	mockService := NewMockService(t)

	expectedOrder := &model.Order{
		UID:         "test123",
		TrackNumber: "TRACK123",
		Entry:       "WBIL",
		Payment:     model.Payment{},
		Delivery:    model.Delivery{},
	}

	mockService.On("GetOrderById", mock.Anything, "test123").
		Return(expectedOrder, nil).
		Once()

	controller := New(mockService)
	handler := controller.GetOrderById()

	req := createTestRequest(t, "test123")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseOrder model.Order
	err := json.Unmarshal(rr.Body.Bytes(), &responseOrder)
	require.NoError(t, err)
	assert.Equal(t, "test123", responseOrder.UID)

	mockService.AssertExpectations(t)
}

func TestGetOrderById_NotFound(t *testing.T) {
	mockService := NewMockService(t)

	mockService.On("GetOrderById", mock.Anything, "nonexistent").
		Return(nil, serviceErrors.ErrNotFound.ForEntity("order")).
		Once()

	controller := New(mockService)
	handler := controller.GetOrderById()

	req := createTestRequest(t, "nonexistent")
	requestRecorder := httptest.NewRecorder()

	handler.ServeHTTP(requestRecorder, req)

	assert.Equal(t, http.StatusNotFound, requestRecorder.Code)
	assert.Equal(t, "order not found\n", requestRecorder.Body.String())

	mockService.AssertExpectations(t)
}

func TestGetOrderById_EmptyOrderUID(t *testing.T) {
	mockService := NewMockService(t)

	controller := New(mockService)
	handler := controller.GetOrderById()

	req := createTestRequest(t, "")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "order_uid is required\n", rr.Body.String())

	mockService.AssertNumberOfCalls(t, "GetOrderById", 0)
}

func createTestRequest(t *testing.T, orderUID string) *http.Request {
	req, err := http.NewRequest("GET", "/order/"+orderUID, nil)
	require.NoError(t, err)

	rctx := chi.NewRouteContext()
	if orderUID != "" {
		rctx.URLParams.Add("order_uid", orderUID)
	}

	return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx))
}
