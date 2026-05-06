package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ILarious/BackForOrder/internal/domain/model"
)

var ErrEmptyUsername = errors.New("order: username is required")

type OrderRepository interface {
	Create(ctx context.Context, username string) (model.Order, error)
	List(ctx context.Context) ([]model.Order, error)
}

type OrderService struct {
	orders OrderRepository
}

func NewOrderService(orders OrderRepository) *OrderService {
	return &OrderService{
		orders: orders,
	}
}

func (s *OrderService) Create(ctx context.Context, username string) (model.Order, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return model.Order{}, ErrEmptyUsername
	}

	return s.orders.Create(ctx, username)
}

func (s *OrderService) List(ctx context.Context) ([]model.Order, error) {
	return s.orders.List(ctx)
}
