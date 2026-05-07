package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ILarious/BackForOrder/internal/domain/model"
)

func TestOutboxDispatcherMarksEventSentAfterPublish(t *testing.T) {
	repo := &outboxRepositoryMock{
		events: []model.OutboxEvent{{ID: 1}},
	}
	publisher := &outboxPublisherMock{}
	dispatcher := NewOutboxDispatcher(repo, publisher, nil)

	dispatcher.dispatch(context.Background())

	if publisher.published != 1 {
		t.Fatalf("published = %d, want 1", publisher.published)
	}
	if repo.sentID != 1 {
		t.Fatalf("sentID = %d, want 1", repo.sentID)
	}
	if repo.failedID != 0 {
		t.Fatalf("failedID = %d, want 0", repo.failedID)
	}
}

func TestOutboxDispatcherMarksEventFailedAfterPublishError(t *testing.T) {
	repo := &outboxRepositoryMock{
		events: []model.OutboxEvent{{ID: 2}},
	}
	publisher := &outboxPublisherMock{err: errors.New("kafka unavailable")}
	dispatcher := NewOutboxDispatcher(repo, publisher, nil)

	dispatcher.dispatch(context.Background())

	if repo.sentID != 0 {
		t.Fatalf("sentID = %d, want 0", repo.sentID)
	}
	if repo.failedID != 2 {
		t.Fatalf("failedID = %d, want 2", repo.failedID)
	}
	if repo.failureReason != "kafka unavailable" {
		t.Fatalf("failure reason = %q, want %q", repo.failureReason, "kafka unavailable")
	}
}

type outboxRepositoryMock struct {
	events        []model.OutboxEvent
	sentID        int64
	failedID      int64
	failureReason string
}

func (r *outboxRepositoryMock) FetchUnsentOutboxEvents(ctx context.Context, limit int) ([]model.OutboxEvent, error) {
	return r.events, nil
}

func (r *outboxRepositoryMock) MarkOutboxEventSent(ctx context.Context, eventID int64) error {
	r.sentID = eventID
	return nil
}

func (r *outboxRepositoryMock) MarkOutboxEventFailed(ctx context.Context, eventID int64, reason string) error {
	r.failedID = eventID
	r.failureReason = reason
	return nil
}

type outboxPublisherMock struct {
	published int
	err       error
}

func (p *outboxPublisherMock) PublishOutboxEvent(ctx context.Context, event model.OutboxEvent) error {
	p.published++
	return p.err
}
