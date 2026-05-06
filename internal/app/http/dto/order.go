package dto

type CreateOrderRequest struct {
	Username string `json:"username"`
}

type OrderResponse struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	FollowersCount int    `json:"followers_count"`
	Status         int    `json:"status"`
}

type ListOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
}
