package main

import (
	"context"
	"log"
	"log/slog"
	"order-service/app/handler"
	"order-service/app/middleware"
	"order-service/app/repository/db"
	stockrepo "order-service/app/repository/stock_repo"
	"order-service/app/usecase"
	"order-service/config"
	"order-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// init logger
	logger.InitLogger()

	// init config
	cfg, err := config.InitConfig(context.Background())
	if err != nil {
		slog.Error("failed to init config", "error", err)
		return
	}

	// init database
	dbConn, err := db.NewPostgres(cfg.Db)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer dbConn.Close()

	reqValidator := validator.New()
	stockRepo := stockrepo.NewStockRepository(cfg.WarehouseService.Host, cfg.InternalAuthHeader)
	orderRepo := db.NewOrderRepository(dbConn)

	orderUsecase := usecase.NewOrderUsecase(orderRepo, stockRepo, cfg)

	orderHandler := handler.NewOrderHandler(orderUsecase, reqValidator)

	// Initialize HTTP web framework
	app := fiber.New()
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(middleware.RequestIDMiddleware())

	handler.SetupRouter(app, orderHandler, cfg)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			slog.Error("Failed to listen", "port", cfg.Port)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Gracefully shutdown")
	err = app.Shutdown()
	if err != nil {
		slog.Warn("Unfortunately the shutdown wasn't smooth", "err", err)
	}
}
