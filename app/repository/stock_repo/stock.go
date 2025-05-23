package stockrepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"order-service/app/domain"
	"order-service/pkg"
	"time"
)

type stockRepository struct {
	httpClient         *http.Client
	baseURL            string
	internalAuthHeader string
}

func NewStockRepository(baseURL string, internalAuthHeader string) domain.StockRepository {
	return &stockRepository{
		httpClient:         &http.Client{Timeout: 30 * time.Second},
		baseURL:            baseURL,
		internalAuthHeader: internalAuthHeader,
	}
}

func (r *stockRepository) CreateReservedStock(ctx context.Context, req domain.ReservedStockCreateRequest) error {
	url := fmt.Sprintf("%s/internal/warehouse-service/reserved-stocks", r.baseURL)
	reqBody, err := json.Marshal(req)
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] CreateReservedStock", "error json Marshal", err)
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] CreateReservedStock", "error http.NewRequestWithContext", err)
		return err
	}

	pkg.AddRequestHeader(ctx, r.internalAuthHeader, httpReq)

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] CreateReservedStock", "error httpClient.Do", err)
		return err
	}
	defer resp.Body.Close()

	var res any
	if err := pkg.DecodeResponseBody(resp, &res); err != nil {
		slog.ErrorContext(ctx, "[stockRepository] CreateReservedStock", "error DecodeResponseBody", err)
		return err
	}

	return nil
}

func (r *stockRepository) UpdateReservedStockStatus(ctx context.Context, orderID int64, req domain.ReservedStockUpdateRequest) error {
	url := fmt.Sprintf("%s/internal/warehouse-service/orders/%d/reserved-stocks/status", r.baseURL, orderID)
	reqBody, err := json.Marshal(req)
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] UpdateReservedStockStatus", "error json Marshal", err)
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(reqBody))
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] UpdateReservedStockStatus", "error http.NewRequestWithContext", err)
		return err
	}

	pkg.AddRequestHeader(ctx, r.internalAuthHeader, httpReq)

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] UpdateReservedStockStatus", "error httpClient.Do", err)
		return err
	}
	defer resp.Body.Close()

	var res any
	if err := pkg.DecodeResponseBody(resp, &res); err != nil {
		slog.ErrorContext(ctx, "[stockRepository] UpdateReservedStockStatus", "error DecodeResponseBody", err)
		return err
	}

	return nil
}
