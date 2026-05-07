package service

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/ILarious/BackForOrder/internal/domain/errors"
	"github.com/ILarious/BackForOrder/internal/domain/model"
)

func TestOrderServiceCreateTrimsUsername(t *testing.T) {
	repo := &orderRepositoryMock{}
	svc := NewOrderService(repo)

	order, err := svc.Create(context.Background(), "  ilarious  ")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if repo.createdUsername != "ilarious" {
		t.Fatalf("created username = %q, want %q", repo.createdUsername, "ilarious")
	}
	if order.Username != "ilarious" {
		t.Fatalf("order username = %q, want %q", order.Username, "ilarious")
	}
}

func TestOrderServiceCreateRejectsEmptyUsername(t *testing.T) {
	svc := NewOrderService(&orderRepositoryMock{})

	_, err := svc.Create(context.Background(), "   ")
	if !errors.Is(err, domainerrors.ErrEmptyUsername) {
		t.Fatalf("Create() error = %v, want %v", err, domainerrors.ErrEmptyUsername)
	}
}

type orderRepositoryMock struct {
	createdUsername string
}

func (r *orderRepositoryMock) Create(ctx context.Context, username string) (model.Order, error) {
	r.createdUsername = username
	return model.Order{ID: 1, Username: username}, nil
}

func (r *orderRepositoryMock) List(ctx context.Context) ([]model.Order, error) {
	return []model.Order{{ID: 1, Username: "ilarious"}}, nil
}

func (r *orderRepositoryMock) UpdateBloggerInfo(ctx context.Context, orderID int64, fullName string, followersCount int, status model.OrderStatus) (model.Order, error) {
	return model.Order{
		ID:             orderID,
		FullName:       fullName,
		FollowersCount: followersCount,
		Status:         status,
	}, nil
}

func (r *orderRepositoryMock) ProcessBloggerInfo(ctx context.Context, messageID, topic string, orderID int64, fullName string, followersCount int, status model.OrderStatus) error {
	return nil
}
