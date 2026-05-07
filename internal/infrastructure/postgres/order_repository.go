package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ILarious/BackForOrder/internal/domain/model"
)

var ErrNilDB = errors.New("order repository: nil db")

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) (*OrderRepository, error) {
	if db == nil {
		return nil, ErrNilDB
	}

	return &OrderRepository{db: db}, nil
}

func (r *OrderRepository) Create(ctx context.Context, username string) (model.Order, error) {
	const query = `
		WITH created_order AS (
			INSERT INTO orders (username)
			VALUES ($1)
			RETURNING id, username, full_name, followers_count, status, created_at, updated_at
		),
		created_event AS (
			INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload)
			SELECT
				'order',
				id,
				'order.created',
				jsonb_build_object('order_id', id, 'username', username)
			FROM created_order
		)
		SELECT id, username, full_name, followers_count, status, created_at, updated_at
		FROM created_order
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Order{}, fmt.Errorf("begin create order: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	order, err := scanOrder(tx.QueryRowContext(ctx, query, username))
	if err != nil {
		return model.Order{}, fmt.Errorf("create order: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return model.Order{}, fmt.Errorf("commit create order: %w", err)
	}

	return order, nil
}

func (r *OrderRepository) List(ctx context.Context) ([]model.Order, error) {
	const query = `
		SELECT id, username, full_name, followers_count, status, created_at, updated_at
		FROM orders
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateBloggerInfo(ctx context.Context, orderID int64, fullName string, followersCount int, status model.OrderStatus) (model.Order, error) {
	const query = `
		UPDATE orders
		SET full_name = $2,
			followers_count = $3,
			status = $4,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, full_name, followers_count, status, created_at, updated_at
	`

	order, err := scanOrder(r.db.QueryRowContext(ctx, query, orderID, fullName, followersCount, int(status)))
	if err != nil {
		return model.Order{}, fmt.Errorf("update order blogger info: %w", err)
	}

	return order, nil
}

func (r *OrderRepository) ProcessBloggerInfo(ctx context.Context, messageID, topic string, orderID int64, fullName string, followersCount int, status model.OrderStatus) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin process blogger info: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	insertResult, err := tx.ExecContext(ctx, `
		INSERT INTO processed_kafka_messages (message_id, topic)
		VALUES ($1, $2)
		ON CONFLICT (message_id) DO NOTHING
	`, messageID, topic)
	if err != nil {
		return fmt.Errorf("save processed kafka message: %w", err)
	}

	inserted, err := insertResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("check processed kafka message: %w", err)
	}
	if inserted == 0 {
		return tx.Commit()
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE orders
		SET full_name = $2,
			followers_count = $3,
			status = $4,
			updated_at = NOW()
		WHERE id = $1
	`, orderID, fullName, followersCount, int(status)); err != nil {
		return fmt.Errorf("update blogger info: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit process blogger info: %w", err)
	}

	return nil
}

func (r *OrderRepository) FetchUnsentOutboxEvents(ctx context.Context, limit int) ([]model.OutboxEvent, error) {
	const query = `
		WITH claimed_events AS (
			SELECT id
			FROM outbox_events
			WHERE sent_at IS NULL
				AND (locked_at IS NULL OR locked_at < NOW() - INTERVAL '5 minutes')
			ORDER BY id
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE outbox_events
		SET locked_at = NOW()
		WHERE id IN (SELECT id FROM claimed_events)
		RETURNING id, aggregate_type, aggregate_id, event_type, payload, created_at
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin fetch unsent outbox events: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	rows, err := tx.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("fetch unsent outbox events: %w", err)
	}
	defer rows.Close()

	events := make([]model.OutboxEvent, 0)
	for rows.Next() {
		var event model.OutboxEvent
		if err := rows.Scan(
			&event.ID,
			&event.AggregateType,
			&event.AggregateID,
			&event.EventType,
			&event.Payload,
			&event.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate outbox events: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit fetch unsent outbox events: %w", err)
	}

	return events, nil
}

func (r *OrderRepository) MarkOutboxEventSent(ctx context.Context, eventID int64) error {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE outbox_events
		SET sent_at = NOW(),
			locked_at = NULL,
			last_error = NULL
		WHERE id = $1
	`, eventID); err != nil {
		return fmt.Errorf("mark outbox event sent: %w", err)
	}

	return nil
}

func (r *OrderRepository) MarkOutboxEventFailed(ctx context.Context, eventID int64, reason string) error {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE outbox_events
		SET attempts = attempts + 1,
			locked_at = NULL,
			last_error = $2
		WHERE id = $1
	`, eventID, reason); err != nil {
		return fmt.Errorf("mark outbox event failed: %w", err)
	}

	return nil
}

type orderScanner interface {
	Scan(dest ...any) error
}

func scanOrder(scanner orderScanner) (model.Order, error) {
	var order model.Order
	var status int

	if err := scanner.Scan(
		&order.ID,
		&order.Username,
		&order.FullName,
		&order.FollowersCount,
		&status,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return model.Order{}, err
	}

	order.Status = model.OrderStatus(status)

	return order, nil
}
