package domain

import (
	"context"
)

type ReservedStockCreateRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
	OrderID   int64 `json:"order_id"`
}

type ReservedStockUpdateRequest struct {
	Status string `json:"status"`
}

type StockRepository interface {
	CreateReservedStock(ctx context.Context, req ReservedStockCreateRequest) error
	UpdateReservedStockStatus(ctx context.Context, orderID int64, req ReservedStockUpdateRequest) error
}
