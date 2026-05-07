package model

import "time"

type OutboxEvent struct {
	ID            int64
	AggregateType string
	AggregateID   int64
	EventType     string
	Payload       []byte
	CreatedAt     time.Time
}
