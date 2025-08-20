// internal/domain/order/service_test.go
package order

import (
	"context"
	"errors"
	"testing"
	"time"

	errors_pkg "wb-L0-task/internal/domain/errors"
	model "wb-L0-task/internal/domain/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOrder_GetOrderById_CacheHit(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	expectedOrder := &model.Order{
		UID:         "test123",
		TrackNumber: "TRACK123",
		Entry:       "WBIL",
	}
	mockCache.On("Get", "test123").Return(*expectedOrder, true).Once()

	result, err := orderService.GetOrderById(context.Background(), "test123")

	require.NoError(t, err)
	assert.Equal(t, expectedOrder, result)
	mockCache.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Exists")
	mockRepo.AssertNotCalled(t, "GetById")
}

func TestOrder_GetOrderById_CacheMiss_DBFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	mockCache.On("Get", "test123").Return(model.Order{}, false).Once()

	mockRepo.On("Exists", mock.Anything, "test123").Return(true, nil).Once()

	expectedOrder := &model.Order{
		UID:         "test123",
		TrackNumber: "TRACK123",
		Entry:       "WBIL",
	}
	mockRepo.On("GetById", mock.Anything, "test123").Return(expectedOrder, nil).Once()

	mockCache.On("Set", "test123", *expectedOrder, time.Duration(0)).Once()

	result, err := orderService.GetOrderById(context.Background(), "test123")

	require.NoError(t, err)
	assert.Equal(t, expectedOrder, result)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestOrder_GetOrderById_CacheMiss_DBNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	mockCache.On("Get", "test123").Return(model.Order{}, false).Once()

	mockRepo.On("Exists", mock.Anything, "test123").Return(false, nil).Once()

	result, err := orderService.GetOrderById(context.Background(), "test123")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, errors_pkg.ErrNotFound))
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetById")
}

func TestOrder_GetOrderById_DBError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	mockCache.On("Get", "test123").Return(model.Order{}, false).Once()

	mockRepo.On("Exists", mock.Anything, "test123").Return(false, assert.AnError).Once()

	result, err := orderService.GetOrderById(context.Background(), "test123")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, assert.AnError, err)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetById")
}

func TestOrder_GetOrderById_DBGetError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	mockCache.On("Get", "test123").Return(model.Order{}, false).Once()

	mockRepo.On("Exists", mock.Anything, "test123").Return(true, nil).Once()
	mockRepo.On("GetById", mock.Anything, "test123").Return(nil, assert.AnError).Once()

	result, err := orderService.GetOrderById(context.Background(), "test123")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, assert.AnError, err)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Set")
}

func TestOrder_InitCache_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	orders := []model.Order{
		{UID: "order1", TrackNumber: "TRACK1"},
		{UID: "order2", TrackNumber: "TRACK2"},
	}

	mockRepo.On("GetOrders", mock.Anything, int32(10)).Return(orders, nil).Once()

	delivery1 := &model.Delivery{Name: "John Doe"}
	delivery2 := &model.Delivery{Name: "Jane Smith"}
	mockRepo.On("GetOrderDelivery", mock.Anything, "order1").Return(delivery1, nil).Once()
	mockRepo.On("GetOrderDelivery", mock.Anything, "order2").Return(delivery2, nil).Once()

	payment1 := &model.Payment{TransactionID: "trans1"}
	payment2 := &model.Payment{TransactionID: "trans2"}
	mockRepo.On("GetOrderPayment", mock.Anything, "order1").Return(payment1, nil).Once()
	mockRepo.On("GetOrderPayment", mock.Anything, "order2").Return(payment2, nil).Once()

	items1 := []model.Item{{ChartID: 1}}
	items2 := []model.Item{{ChartID: 2}}
	mockRepo.On("GetOrderItems", mock.Anything, "order1").Return(items1, nil).Once()
	mockRepo.On("GetOrderItems", mock.Anything, "order2").Return(items2, nil).Once()

	expectedOrder1 := orders[0]
	expectedOrder1.Delivery = *delivery1
	expectedOrder1.Payment = *payment1
	expectedOrder1.Items = items1

	expectedOrder2 := orders[1]
	expectedOrder2.Delivery = *delivery2
	expectedOrder2.Payment = *payment2
	expectedOrder2.Items = items2

	mockCache.On("Set", "order1", expectedOrder1, 30*time.Second).Once()
	mockCache.On("Set", "order2", expectedOrder2, 30*time.Second).Once()

	err := orderService.InitCache(context.Background())

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestOrder_InitCache_GetOrdersError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	mockRepo.On("GetOrders", mock.Anything, int32(10)).Return(nil, assert.AnError).Once()

	err := orderService.InitCache(context.Background())

	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Set")
}

func TestOrder_InitCache_DeliveryError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	orders := []model.Order{
		{UID: "order1", TrackNumber: "TRACK1"},
	}

	mockRepo.On("GetOrders", mock.Anything, int32(10)).Return(orders, nil).Once()
	mockRepo.On("GetOrderDelivery", mock.Anything, "order1").Return(nil, assert.AnError).Once()

	err := orderService.InitCache(context.Background())

	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Set")
}

func TestOrder_InitCache_PaymentError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	orders := []model.Order{
		{UID: "order1", TrackNumber: "TRACK1"},
	}

	mockRepo.On("GetOrders", mock.Anything, int32(10)).Return(orders, nil).Once()
	mockRepo.On("GetOrderDelivery", mock.Anything, "order1").Return(&model.Delivery{}, nil).Once()
	mockRepo.On("GetOrderPayment", mock.Anything, "order1").Return(nil, assert.AnError).Once()

	err := orderService.InitCache(context.Background())

	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Set")
}

func TestOrder_InitCache_ItemsError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache[model.Order])

	orderService := New(mockCache, mockRepo)

	orders := []model.Order{
		{UID: "order1", TrackNumber: "TRACK1"},
	}

	mockRepo.On("GetOrders", mock.Anything, int32(10)).Return(orders, nil).Once()
	mockRepo.On("GetOrderDelivery", mock.Anything, "order1").Return(&model.Delivery{}, nil).Once()
	mockRepo.On("GetOrderPayment", mock.Anything, "order1").Return(&model.Payment{}, nil).Once()
	mockRepo.On("GetOrderItems", mock.Anything, "order1").Return(nil, assert.AnError).Once()

	err := orderService.InitCache(context.Background())

	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Set")
}
