package handler

import (
	"order-service/app/middleware"
	"order-service/config"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, orderHandler *OrderHandler, cfg *config.Config) {
	// Setup routes
	apiGroup := app.Group("/order-service").Use(middleware.Auth(cfg.Jwt.SecretKey))
	callback := app.Group("/callback/order-service").Use(middleware.AuthPayment(cfg))

	apiGroup.Get("/orders/:id", orderHandler.GetOrderByID)
	apiGroup.Get("/orders", orderHandler.GetListByUserID)
	apiGroup.Post("/orders", orderHandler.CreateOrder)

	// callback payment update order status
	callback.Post("/orders", orderHandler.UpdateStatusOrder)
}
