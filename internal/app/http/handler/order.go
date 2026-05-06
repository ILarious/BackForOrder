package handler

type OrderService interface {
}

type OrderHandler struct {
	Orders OrderService
}

func NewOrderHandler(orders OrderService) *OrderHandler {
	return &OrderHandler{Orders: orders}
}
