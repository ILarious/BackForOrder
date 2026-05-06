package model

import "time"

type OrderStatus int

const (
	OrderStatusNew OrderStatus = iota
	OrderStatusProcessing
	OrderStatusDone
	OrderStatusFailed
)

type Order struct {
	ID             int64
	Username       string
	FullName       string
	FollowersCount int
	Status         OrderStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
