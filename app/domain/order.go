package domain

import (
	"context"
	"database/sql"
	"time"
)

type OrderStatus string

const (
	OrderStatusWaitingPayment OrderStatus = "waiting_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusCancelled      OrderStatus = "cancelled"
)

type Order struct {
	ID        int64       `json:"id"`
	ProductID int64       `json:"product_id"`
	Quantity  int64       `json:"quantity"`
	UserID    int64       `json:"user_id"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderUpdateStatusRequest struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status" validate:"required,oneof=paid cancelled"`
}

type OrderCreateRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order, tx *sql.Tx) error
	GetOrderByID(ctx context.Context, id int64) (Order, error)
	UpdateStatusOrder(ctx context.Context, id int64, status string, tx *sql.Tx) error
	GetListByUserID(ctx context.Context, userID int64) ([]Order, error)

	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}

type OrderUsecase interface {
	CreateOrder(ctx context.Context, userID int64, req OrderCreateRequest) (Order, error)
	UpdateStatusOrder(ctx context.Context, req OrderUpdateStatusRequest) error
	GetListByUserID(ctx context.Context, userID int64) ([]Order, error)
	GetOrderByID(ctx context.Context, userID int64, id int64) (Order, error)
}
