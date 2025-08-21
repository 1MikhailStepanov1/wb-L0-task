package order

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	serviceErrors "wb-L0-task/internal/domain/errors"
	models "wb-L0-task/internal/domain/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestKafkaConsumerService_SaveOrder_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	validOrder := &models.Order{
		UID: "test123",
		Delivery: models.Delivery{
			Phone: "+79161234567",
			Email: "test@example.com",
		},
		Payment: models.Payment{
			PaymentDT:    time.Now().Truncate(time.Second),
			GoodsTotal:   1000,
			DeliveryCost: 500,
			CustomFee:    100,
			Amount:       1600,
		},
		Items: []models.Item{
			{
				Price:      1000,
				Sale:       0,
				TotalPrice: 1000,
			},
		},
	}

	orderJSON, err := json.Marshal(validOrder)
	require.NoError(t, err)

	mockRepo.On("Save", mock.Anything, validOrder).Return(nil).Once()

	err = service.SaveOrder(context.Background(), orderJSON)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestKafkaConsumerService_SaveOrder_EmptyOrderUID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	order := &models.Order{
		UID: "",
		Delivery: models.Delivery{
			Phone: "+79161234567",
			Email: "test@example.com",
		},
	}

	orderJSON, err := json.Marshal(order)
	require.NoError(t, err)

	err = service.SaveOrder(context.Background(), orderJSON)

	require.Error(t, err)
	assert.True(t, errors.Is(err, serviceErrors.ErrInvalidEntity))
	mockRepo.AssertNotCalled(t, "Save")
}

func TestKafkaConsumerService_SaveOrder_InvalidPhoneNumber(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	testCases := []struct {
		name  string
		phone string
	}{
		{"empty phone", ""},
		{"too short", "123"},
		{"invalid format", "invalid-phone"},
		{"too long", "+7916123456789012345"},
		{"letters", "abc123def"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := &models.Order{
				UID: "test123",
				Delivery: models.Delivery{
					Phone: tc.phone,
					Email: "test@example.com",
				},
			}

			orderJSON, err := json.Marshal(order)
			require.NoError(t, err)

			err = service.SaveOrder(context.Background(), orderJSON)

			require.Error(t, err)
			assert.True(t, errors.Is(err, serviceErrors.ErrInvalidEntity))
			mockRepo.AssertNotCalled(t, "Save")
		})
	}
}

func TestKafkaConsumerService_SaveOrder_ValidPhoneNumbers(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	validPhones := []string{
		"+79161234567",
		"89161234567",
		"79161234567",
		"9161234567",
		"+123456789012345",
	}

	for _, phone := range validPhones {
		t.Run("phone_"+phone, func(t *testing.T) {
			order := &models.Order{
				UID: "test123",
				Delivery: models.Delivery{
					Phone: phone,
					Email: "test@example.com",
				},
				Payment: models.Payment{
					PaymentDT:    time.Now().Truncate(time.Second),
					GoodsTotal:   1000,
					DeliveryCost: 500,
					CustomFee:    100,
					Amount:       1600,
				},
				Items: []models.Item{
					{
						Price:      1000,
						Sale:       0,
						TotalPrice: 1000,
					},
				},
			}

			orderJSON, err := json.Marshal(order)
			require.NoError(t, err)

			mockRepo.On("Save", mock.Anything, order).Return(nil).Once()

			err = service.SaveOrder(context.Background(), orderJSON)

			require.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestKafkaConsumerService_SaveOrder_InvalidEmail(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	testCases := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"no @", "invalid-email"},
		{"no domain", "test@"},
		{"no username", "@example.com"},
		{"invalid domain", "test@.com"},
		{"multiple @", "test@@example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := &models.Order{
				UID: "test123",
				Delivery: models.Delivery{
					Phone: "+79161234567",
					Email: tc.email,
				},
			}

			orderJSON, err := json.Marshal(order)
			require.NoError(t, err)

			err = service.SaveOrder(context.Background(), orderJSON)

			require.Error(t, err)
			assert.True(t, errors.Is(err, serviceErrors.ErrInvalidEntity))
			mockRepo.AssertNotCalled(t, "Save")
		})
	}
}

func TestKafkaConsumerService_SaveOrder_ValidEmails(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	validEmails := []string{
		"test@example.com",
		"user.name@example.com",
		"user.name+tag@example.com",
		"user@sub.domain.com",
		"user@example.co.uk",
		"123@example.com",
	}

	for _, email := range validEmails {
		t.Run("email_"+email, func(t *testing.T) {
			order := &models.Order{
				UID: "test123",
				Delivery: models.Delivery{
					Phone: "+79161234567",
					Email: email,
				},
				Payment: models.Payment{
					PaymentDT:    time.Now().Truncate(time.Second),
					GoodsTotal:   1000,
					DeliveryCost: 500,
					CustomFee:    100,
					Amount:       1600,
				},
				Items: []models.Item{
					{
						Price:      1000,
						Sale:       0,
						TotalPrice: 1000,
					},
				},
			}

			orderJSON, err := json.Marshal(order)
			require.NoError(t, err)

			mockRepo.On("Save", mock.Anything, order).Return(nil).Once()

			err = service.SaveOrder(context.Background(), orderJSON)

			require.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestKafkaConsumerService_SaveOrder_InvalidItemTotalPrice(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	testCases := []struct {
		name  string
		item  models.Item
		error string
	}{
		{
			"sale calculation error",
			models.Item{Price: 1000, Sale: 10, TotalPrice: 1000},
			"order.item.total_price",
		},
		{
			"zero price with sale",
			models.Item{Price: 0, Sale: 10, TotalPrice: 0},
			"order.item.total_price",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := &models.Order{
				UID: "test123",
				Delivery: models.Delivery{
					Phone: "+79161234567",
					Email: "test@example.com",
				},
				Payment: models.Payment{
					PaymentDT:    time.Now().Truncate(time.Second),
					GoodsTotal:   1000,
					DeliveryCost: 500,
					CustomFee:    100,
					Amount:       1600,
				},
				Items: []models.Item{tc.item},
			}

			orderJSON, err := json.Marshal(order)
			require.NoError(t, err)

			err = service.SaveOrder(context.Background(), orderJSON)

			require.Error(t, err)
			assert.True(t, errors.Is(err, serviceErrors.ErrInvalidEntity))
			mockRepo.AssertNotCalled(t, "Save")
		})
	}
}

func TestKafkaConsumerService_SaveOrder_InvalidGoodsTotal(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	order := &models.Order{
		UID: "test123",
		Delivery: models.Delivery{
			Phone: "+79161234567",
			Email: "test@example.com",
		},
		Payment: models.Payment{
			PaymentDT:    time.Now().Truncate(time.Second),
			GoodsTotal:   2000,
			DeliveryCost: 500,
			CustomFee:    100,
			Amount:       2600,
		},
		Items: []models.Item{
			{
				Price:      1000,
				Sale:       0,
				TotalPrice: 1000,
			},
		},
	}

	orderJSON, err := json.Marshal(order)
	require.NoError(t, err)

	err = service.SaveOrder(context.Background(), orderJSON)

	require.Error(t, err)
	assert.True(t, errors.Is(err, serviceErrors.ErrInvalidEntity))
	mockRepo.AssertNotCalled(t, "Save")
}

func TestKafkaConsumerService_SaveOrder_InvalidPaymentAmount(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	order := &models.Order{
		UID: "test123",
		Delivery: models.Delivery{
			Phone: "+79161234567",
			Email: "test@example.com",
		},
		Payment: models.Payment{
			PaymentDT:    time.Now().Truncate(time.Second),
			GoodsTotal:   1000,
			DeliveryCost: 500,
			CustomFee:    100,
			Amount:       2000,
		},
		Items: []models.Item{
			{
				Price:      1000,
				Sale:       0,
				TotalPrice: 1000,
			},
		},
	}

	orderJSON, err := json.Marshal(order)
	require.NoError(t, err)

	err = service.SaveOrder(context.Background(), orderJSON)

	require.Error(t, err)
	assert.True(t, errors.Is(err, serviceErrors.ErrInvalidEntity))
	mockRepo.AssertNotCalled(t, "Save")
}

func TestKafkaConsumerService_SaveOrder_MultipleItems(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	order := &models.Order{
		UID: "test123",
		Delivery: models.Delivery{
			Phone: "+79161234567",
			Email: "test@example.com",
		},
		Payment: models.Payment{
			PaymentDT:    time.Now().Truncate(time.Second),
			GoodsTotal:   2250,
			DeliveryCost: 500,
			CustomFee:    100,
			Amount:       2850,
		},
		Items: []models.Item{
			{
				Price:      1000,
				Sale:       0,
				TotalPrice: 1000,
			},
			{
				Price:      1500,
				Sale:       50,
				TotalPrice: 750,
			},
			{
				Price:      1000,
				Sale:       50,
				TotalPrice: 500,
			},
		},
	}

	orderJSON, err := json.Marshal(order)
	require.NoError(t, err)

	mockRepo.On("Save", mock.Anything, order).Return(nil).Once()

	err = service.SaveOrder(context.Background(), orderJSON)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestKafkaConsumerService_isValidOrder_EmptyOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewKafkaConsumerService(mockRepo)

	emptyOrder := &models.Order{}

	err := service.isValidOrder(emptyOrder)
	require.Error(t, err)
	assert.True(t, errors.Is(err, serviceErrors.ErrInvalidEntity))
}
