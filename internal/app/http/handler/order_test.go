package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domainerrors "github.com/ILarious/BackForOrder/internal/domain/errors"
	"github.com/ILarious/BackForOrder/internal/domain/model"
)

func TestOrderHandlerCreate(t *testing.T) {
	handler := NewOrderHandler(&orderServiceMock{
		createOrder: model.Order{ID: 1, Username: "ilarious"},
	})
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"username":"ilarious"}`))
	response := httptest.NewRecorder()

	handler.Create(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusCreated)
	}

	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["username"] != "ilarious" {
		t.Fatalf("username = %v, want ilarious", payload["username"])
	}
}

func TestOrderHandlerCreateRejectsInvalidJSON(t *testing.T) {
	handler := NewOrderHandler(&orderServiceMock{})
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{`))
	response := httptest.NewRecorder()

	handler.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestOrderHandlerCreateMapsDomainErrorToBadRequest(t *testing.T) {
	handler := NewOrderHandler(&orderServiceMock{createErr: domainerrors.ErrEmptyUsername})
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"username":""}`))
	response := httptest.NewRecorder()

	handler.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestOrderHandlerList(t *testing.T) {
	handler := NewOrderHandler(&orderServiceMock{
		listOrders: []model.Order{{ID: 1, Username: "ilarious"}},
	})
	request := httptest.NewRequest(http.MethodGet, "/orders", nil)
	response := httptest.NewRecorder()

	handler.List(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), `"orders"`) {
		t.Fatalf("response body = %s, want orders field", response.Body.String())
	}
}

func TestOrderHandlerListReturnsInternalServerError(t *testing.T) {
	handler := NewOrderHandler(&orderServiceMock{listErr: errors.New("db down")})
	request := httptest.NewRequest(http.MethodGet, "/orders", nil)
	response := httptest.NewRecorder()

	handler.List(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}

type orderServiceMock struct {
	createOrder model.Order
	createErr   error
	listOrders  []model.Order
	listErr     error
}

func (s *orderServiceMock) Create(ctx context.Context, username string) (model.Order, error) {
	return s.createOrder, s.createErr
}

func (s *orderServiceMock) List(ctx context.Context) ([]model.Order, error) {
	return s.listOrders, s.listErr
}
