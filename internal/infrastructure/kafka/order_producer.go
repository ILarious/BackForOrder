package kafka

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ILarious/BackForOrder/internal/domain/model"
	"github.com/segmentio/kafka-go"
)

type OrderProducer struct {
	writer *kafka.Writer
}

func NewOrderProducer(brokers []string, topic string) *OrderProducer {
	return &OrderProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Async:        false,
		},
	}
}

func (p *OrderProducer) PublishOutboxEvent(ctx context.Context, event model.OutboxEvent) error {
	var request OrderRequest
	if err := json.Unmarshal(event.Payload, &request); err != nil {
		return err
	}
	request.EventID = event.ID

	payload, err := json.Marshal(OrderRequest{
		EventID:  request.EventID,
		OrderID:  request.OrderID,
		Username: request.Username,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(event.ID, 10)),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event_id", Value: []byte(strconv.FormatInt(event.ID, 10))},
			{Key: "event_type", Value: []byte(event.EventType)},
		},
	})
}

func (p *OrderProducer) Close() error {
	return p.writer.Close()
}
