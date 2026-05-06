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
		INSERT INTO orders (username)
		VALUES ($1)
		RETURNING id, username, full_name, followers_count, status, created_at, updated_at
	`

	order, err := scanOrder(r.db.QueryRowContext(ctx, query, username))
	if err != nil {
		return model.Order{}, fmt.Errorf("create order: %w", err)
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
