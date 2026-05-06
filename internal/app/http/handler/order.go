package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ILarious/BackForOrder/internal/app/http/dto"
	"github.com/ILarious/BackForOrder/internal/domain/model"
	"github.com/ILarious/BackForOrder/internal/domain/service"
)

type OrderService interface {
	Create(ctx context.Context, username string) (model.Order, error)
	List(ctx context.Context) ([]model.Order, error)
}

type OrderHandler struct {
	Orders OrderService
}

func NewOrderHandler(orders OrderService) *OrderHandler {
	return &OrderHandler{Orders: orders}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	order, err := h.Orders.Create(r.Context(), request.Username)
	if err != nil {
		if errors.Is(err, service.ErrEmptyUsername) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create order"})
		return
	}

	writeJSON(w, http.StatusCreated, orderResponse(order))
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	orders, err := h.Orders.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list orders"})
		return
	}

	response := dto.ListOrdersResponse{
		Orders: make([]dto.OrderResponse, 0, len(orders)),
	}
	for _, order := range orders {
		response.Orders = append(response.Orders, orderResponse(order))
	}

	writeJSON(w, http.StatusOK, response)
}

func orderResponse(order model.Order) dto.OrderResponse {
	return dto.OrderResponse{
		ID:             int(order.ID),
		Username:       order.Username,
		FullName:       order.FullName,
		FollowersCount: order.FollowersCount,
		Status:         int(order.Status),
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(payload)
}
