package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ILarious/BackForOrder/internal/domain/model"
	"github.com/ILarious/BackForOrder/pkg/worker_pool"
	"github.com/segmentio/kafka-go"
)

type OrderResultHandler interface {
	ProcessBloggerInfo(ctx context.Context, messageID, topic string, orderID int64, fullName string, followersCount int, status model.OrderStatus) error
}

type OrderConsumer struct {
	reader  *kafka.Reader
	handler OrderResultHandler
	workers *worker_pool.WorkerPool
}

func NewOrderConsumer(brokers []string, topic, groupID string, handler OrderResultHandler, workers *worker_pool.WorkerPool) *OrderConsumer {
	return &OrderConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		handler: handler,
		workers: workers,
	}
}

func (c *OrderConsumer) Run(ctx context.Context) {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			log.Printf("failed to fetch kafka order result: %v", err)
			continue
		}

		if err := c.submit(ctx, message); err != nil {
			log.Printf("failed to submit kafka order result task: %v", err)
		}
	}
}

func (c *OrderConsumer) submit(ctx context.Context, message kafka.Message) error {
	if c.workers == nil {
		return c.handleMessage(ctx, message)
	}

	return c.workers.Submit(func() {
		if err := c.handleMessage(ctx, message); err != nil {
			log.Printf("failed to handle kafka order result: %v", err)
		}
	})
}

func (c *OrderConsumer) handleMessage(ctx context.Context, message kafka.Message) error {
	var response OrderResponse
	if err := json.Unmarshal(message.Value, &response); err != nil {
		return err
	}

	messageID := response.EventID
	if messageID == "" {
		messageID = string(message.Key)
	}
	if messageID == "" {
		messageID = fmt.Sprintf("%s:%d:%d", message.Topic, message.Partition, message.Offset)
	}

	if err := c.handler.ProcessBloggerInfo(
		ctx,
		messageID,
		message.Topic,
		response.OrderID,
		response.FullName,
		response.FollowersCount,
		model.OrderStatus(response.Status),
	); err != nil {
		return err
	}

	return c.reader.CommitMessages(ctx, message)
}

func (c *OrderConsumer) Close() error {
	return c.reader.Close()
}
