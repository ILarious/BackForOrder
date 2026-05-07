package service

import (
	"context"
	"strings"

	domainerrors "github.com/ILarious/BackForOrder/internal/domain/errors"
	"github.com/ILarious/BackForOrder/internal/domain/model"
)

type OrderRepository interface {
	Create(ctx context.Context, username string) (model.Order, error)
	List(ctx context.Context) ([]model.Order, error)
	UpdateBloggerInfo(ctx context.Context, orderID int64, fullName string, followersCount int, status model.OrderStatus) (model.Order, error)
	ProcessBloggerInfo(ctx context.Context, messageID, topic string, orderID int64, fullName string, followersCount int, status model.OrderStatus) error
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
		return model.Order{}, domainerrors.ErrEmptyUsername
	}

	return s.orders.Create(ctx, username)
}

func (s *OrderService) List(ctx context.Context) ([]model.Order, error) {
	return s.orders.List(ctx)
}

func (s *OrderService) UpdateBloggerInfo(ctx context.Context, orderID int64, fullName string, followersCount int, status model.OrderStatus) (model.Order, error) {
	return s.orders.UpdateBloggerInfo(ctx, orderID, fullName, followersCount, status)
}

func (s *OrderService) ProcessBloggerInfo(ctx context.Context, messageID, topic string, orderID int64, fullName string, followersCount int, status model.OrderStatus) error {
	return s.orders.ProcessBloggerInfo(ctx, messageID, topic, orderID, fullName, followersCount, status)
}
