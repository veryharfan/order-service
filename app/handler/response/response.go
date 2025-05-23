package response

import (
	"errors"
	"order-service/app/domain"

	"github.com/gofiber/fiber/v2"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Success[T any](data T) *Response[T] {
	return &Response[T]{
		Success: true,
		Data:    data,
	}
}

func Error(err error) *Response[any] {
	return &Response[any]{
		Success: false,
		Error:   err.Error(),
	}
}

func FromError(err error) (int, *Response[any]) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return fiber.StatusBadRequest, Error(err)
	case errors.Is(err, domain.ErrInvalidRequest):
		return fiber.StatusBadRequest, Error(err)
	case errors.Is(err, domain.ErrUnauthorized):
		return fiber.StatusUnauthorized, Error(err)
	case errors.Is(err, domain.ErrNotFound):
		return fiber.StatusNotFound, Error(err)
	case errors.Is(err, domain.ErrBadRequest):
		return fiber.StatusBadRequest, Error(err)
	default:
		return fiber.StatusInternalServerError, Error(domain.ErrInternal)
	}
}
