package db

import (
	"context"
	"database/sql"
	"log/slog"
	"order-service/app/domain"
	"time"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order, tx *sql.Tx) error {
	query := `INSERT INTO orders (product_id, quantity, user_id, status, created_at, updated_at, expired_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := tx.QueryRowContext(ctx, query,
		order.ProductID,
		order.Quantity,
		order.UserID,
		order.Status,
		time.Now(),
		time.Now(),
		order.ExpiredAt,
	).Scan(&order.ID)
	if err != nil {
		slog.ErrorContext(ctx, "[orderRepository] CreateOrder", "failed to create order", err)
		return err
	}
	return nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id int64) (domain.Order, error) {
	query := `SELECT id, product_id, quantity, user_id, status, created_at, updated_at, expired_at
		FROM orders WHERE id = $1`
	order := domain.Order{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.ProductID,
		&order.Quantity,
		&order.UserID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.ExpiredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "[orderRepository] GetOrderByID", "order not found", err)
			return order, domain.ErrNotFound
		}
		slog.ErrorContext(ctx, "[orderRepository] GetOrderByID", "failed to get order by ID", err)
		return order, err
	}
	return order, nil
}

func (r *orderRepository) UpdateStatusOrder(ctx context.Context, id int64, status string, tx *sql.Tx) error {
	query := `UPDATE orders SET status = $1, updated_at = now() WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, status, id)
	if err != nil {
		slog.ErrorContext(ctx, "[orderRepository] UpdateStatusOrder", "failed to update order status", err)
		return err
	}
	return nil
}

func (r *orderRepository) GetListByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	query := `SELECT id, product_id, quantity, user_id, status, created_at, updated_at, expired_at
		FROM orders WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		slog.ErrorContext(ctx, "[orderRepository] GetListByUserID", "failed to get orders by user ID", err)
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		order := domain.Order{}
		err := rows.Scan(
			&order.ID,
			&order.ProductID,
			&order.Quantity,
			&order.UserID,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.ExpiredAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, "[orderRepository] GetListByUserID", "scan error", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *orderRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "[orderRepository] WithTransaction", "failed to begin transaction", err)
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				slog.ErrorContext(ctx, "[orderRepository] WithTransaction", "failed to rollback transaction", err)
			}
		}
	}()

	if err := fn(ctx, tx); err != nil {
		slog.ErrorContext(ctx, "[orderRepository] WithTransaction", "function error", err)
		return err
	}
	if err := tx.Commit(); err != nil {
		slog.ErrorContext(ctx, "[orderRepository] WithTransaction", "failed to commit transaction", err)
		return err
	}
	return nil
}

func (r *orderRepository) GetExpiredOrders(ctx context.Context) ([]domain.Order, error) {
	query := `SELECT id, product_id, quantity, user_id, status, created_at, updated_at, expired_at
		FROM orders WHERE status = 'waiting_payment' AND expired_at < now()`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.ErrorContext(ctx, "[orderRepository] GetExpiredOrders", "failed to get expired orders", err)
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		order := domain.Order{}
		err := rows.Scan(
			&order.ID,
			&order.ProductID,
			&order.Quantity,
			&order.UserID,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.ExpiredAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, "[orderRepository] GetExpiredOrders", "scan error", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
