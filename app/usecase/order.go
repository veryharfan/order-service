package usecase

import (
	"context"
	"database/sql"
	"log/slog"
	"order-service/app/domain"
	"order-service/config"
)

type orderUsecase struct {
	orderRepository domain.OrderRepository
	stockRepository domain.StockRepository
	cfg             *config.Config
}

func NewOrderUsecase(orderRepository domain.OrderRepository, stockRepository domain.StockRepository, cfg *config.Config) domain.OrderUsecase {
	return &orderUsecase{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
		cfg:             cfg,
	}
}

func (u *orderUsecase) CreateOrder(ctx context.Context, userID int64, req domain.OrderCreateRequest) (domain.Order, error) {
	order := domain.Order{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UserID:    userID,
		Status:    domain.OrderStatusWaitingPayment,
	}

	err := u.orderRepository.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := u.orderRepository.CreateOrder(ctx, &order, tx)
		if err != nil {
			slog.ErrorContext(ctx, "[orderUsecase] CreateOrder", "failed to create order", err)
			return err
		}

		reservedStockReq := domain.ReservedStockCreateRequest{
			ProductID: order.ProductID,
			Quantity:  order.Quantity,
			OrderID:   order.ID,
		}

		if err := u.stockRepository.CreateReservedStock(ctx, reservedStockReq); err != nil {
			slog.ErrorContext(ctx, "[orderUsecase] CreateOrder", "failed to create reserved stock", err)
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (u *orderUsecase) GetListByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	orders, err := u.orderRepository.GetListByUserID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "[orderUsecase] GetListByUserID", "failed to get list by user ID", err)
		return nil, err
	}
	return orders, nil
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, userID, id int64) (domain.Order, error) {
	order, err := u.orderRepository.GetOrderByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "[orderUsecase] GetOrderByID", "failed to get order by ID", err)
		return domain.Order{}, err
	}

	if order.UserID != userID {
		slog.ErrorContext(ctx, "[orderUsecase] GetOrderByID", "user ID mismatch", "unauthorized access")
		return domain.Order{}, domain.ErrUnauthorized
	}

	return order, nil
}

func (u *orderUsecase) UpdateStatusOrder(ctx context.Context, req domain.OrderUpdateStatusRequest) error {
	var reservedStockReq domain.ReservedStockUpdateRequest
	if req.Status == string(domain.OrderStatusCancelled) {
		reservedStockReq.Status = "cancelled"
	} else if req.Status == string(domain.OrderStatusPaid) {
		reservedStockReq.Status = "completed"
	} else {
		slog.ErrorContext(ctx, "[orderUsecase] UpdateStatusOrder", "invalid status", req.Status)
		return domain.ErrBadRequest
	}

	if err := u.orderRepository.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := u.orderRepository.UpdateStatusOrder(ctx, req.OrderID, req.Status, tx)
		if err != nil {
			slog.ErrorContext(ctx, "[orderUsecase] UpdateStatusOrder", "failed to update order status", err)
			return err
		}

		if err := u.stockRepository.UpdateReservedStockStatus(ctx, req.OrderID, reservedStockReq); err != nil {
			slog.ErrorContext(ctx, "[orderUsecase] UpdateStatusOrder", "failed to update reserved stock status", err)
			return err
		}

		return nil
	}); err != nil {
		slog.ErrorContext(ctx, "[orderUsecase] UpdateStatusOrder", "transaction", err)
		return err
	}

	slog.InfoContext(ctx, "[orderUsecase] success UpdateStatusOrder", "order_id", req)
	return nil
}
