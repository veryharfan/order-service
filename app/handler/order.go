package handler

import (
	"log/slog"
	"order-service/app/domain"
	"order-service/app/handler/response"
	"order-service/pkg/ctxutil"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	OrderUsecase domain.OrderUsecase
	validator    *validator.Validate
}

func NewOrderHandler(orderUsecase domain.OrderUsecase, validator *validator.Validate) *OrderHandler {
	return &OrderHandler{
		OrderUsecase: orderUsecase,
		validator:    validator,
	}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var order domain.OrderCreateRequest
	if err := c.BodyParser(&order); err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] CreateOrder", "body", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(order); err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] CreateOrder", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	userID, err := ctxutil.GetUserIDCtx(c.Context())
	if err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] CreateOrder", "getUserIDCtx", err)
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
	}

	res, err := h.OrderUsecase.CreateOrder(c.Context(), userID, order)
	if err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] CreateOrder", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res))
}

func (h *OrderHandler) GetListByUserID(c *fiber.Ctx) error {
	userID, err := ctxutil.GetUserIDCtx(c.Context())
	if err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] GetListByUserID", "getUserIDCtx", err)
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
	}

	res, err := h.OrderUsecase.GetListByUserID(c.Context(), userID)
	if err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] GetListByUserID", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res))
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	idstr := c.Params("id")
	if idstr == "" {
		slog.ErrorContext(c.Context(), "[OrderHandler] GetOrderByID", "params", "order ID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil || id <= 0 {
		slog.ErrorContext(c.Context(), "[OrderHandler] GetOrderByID", "params:"+idstr, err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}
	userID, err := ctxutil.GetUserIDCtx(c.Context())
	if err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] GetOrderByID", "getUserIDCtx", err)
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
	}
	res, err := h.OrderUsecase.GetOrderByID(c.Context(), userID, id)
	if err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] GetOrderByID", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}
	return c.Status(fiber.StatusOK).JSON(response.Success(res))
}

func (h *OrderHandler) UpdateStatusOrder(c *fiber.Ctx) error {
	var req domain.OrderUpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] UpdateStatusOrder", "body", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] UpdateStatusOrder", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if err := h.OrderUsecase.UpdateStatusOrder(c.Context(), req); err != nil {
		slog.ErrorContext(c.Context(), "[OrderHandler] UpdateStatusOrder", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any](nil))
}
