package service

import (
	"context"
	"log"
	"time"

	"github.com/ILarious/BackForOrder/internal/domain/model"
	"github.com/ILarious/BackForOrder/pkg/worker_pool"
)

type OutboxRepository interface {
	FetchUnsentOutboxEvents(ctx context.Context, limit int) ([]model.OutboxEvent, error)
	MarkOutboxEventSent(ctx context.Context, eventID int64) error
	MarkOutboxEventFailed(ctx context.Context, eventID int64, reason string) error
}

type OutboxPublisher interface {
	PublishOutboxEvent(ctx context.Context, event model.OutboxEvent) error
}

type OutboxDispatcher struct {
	events    OutboxRepository
	publisher OutboxPublisher
	workers   *worker_pool.WorkerPool
	interval  time.Duration
	batchSize int
}

func NewOutboxDispatcher(events OutboxRepository, publisher OutboxPublisher, workers *worker_pool.WorkerPool) *OutboxDispatcher {
	return &OutboxDispatcher{
		events:    events,
		publisher: publisher,
		workers:   workers,
		interval:  2 * time.Second,
		batchSize: 50,
	}
}

func (d *OutboxDispatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	d.dispatch(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.dispatch(ctx)
		}
	}
}

func (d *OutboxDispatcher) dispatch(ctx context.Context) {
	events, err := d.events.FetchUnsentOutboxEvents(ctx, d.batchSize)
	if err != nil {
		log.Printf("failed to fetch outbox events: %v", err)
		return
	}

	for _, event := range events {
		event := event
		if d.workers == nil {
			d.publish(ctx, event)
			continue
		}

		if err := d.workers.Submit(func() {
			d.publish(context.Background(), event)
		}); err != nil {
			log.Printf("failed to submit outbox event %d: %v", event.ID, err)
		}
	}
}

func (d *OutboxDispatcher) publish(ctx context.Context, event model.OutboxEvent) {
	if err := d.publisher.PublishOutboxEvent(ctx, event); err != nil {
		log.Printf("failed to publish outbox event %d: %v", event.ID, err)
		if markErr := d.events.MarkOutboxEventFailed(context.Background(), event.ID, err.Error()); markErr != nil {
			log.Printf("failed to mark outbox event %d failed: %v", event.ID, markErr)
		}
		return
	}

	if err := d.events.MarkOutboxEventSent(context.Background(), event.ID); err != nil {
		log.Printf("failed to mark outbox event %d sent: %v", event.ID, err)
	}
}
