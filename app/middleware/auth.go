package middleware

import (
	"log/slog"
	"order-service/app/domain"
	"order-service/app/handler/response"
	"order-service/config"
	"order-service/pkg"
	"order-service/pkg/ctxutil"

	"github.com/gofiber/fiber/v2"
)

type AuthHeader string

const (
	AuthPaymentHeaderKey AuthHeader = "X-Payment-Auth"
)

func AuthPayment(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the auth header from the request
		authHeader := c.Get(string(AuthPaymentHeaderKey))
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
		}
		// Check if the auth header is valid (you can implement your own logic here)
		if authHeader != cfg.AuthPaymentHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
		}

		return c.Next()
	}
}

func Auth(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, err := pkg.GetTokenFromHeaders(c.Get("Authorization"))
		if err != nil {
			slog.ErrorContext(c.Context(), "[middleware] Auth", "GetTokenFromHeaders", err)
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
		}

		claims, err := pkg.ParseJwtToken(token, secretKey)
		if err != nil {
			slog.ErrorContext(c.Context(), "[middleware] Auth", "ParseJwtToken", err)
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
		}

		if claims.UID == 0 {
			slog.ErrorContext(c.Context(), "[middleware] Auth", "userID", "0")
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
		}

		c.Locals(ctxutil.UserIDKey, claims.UID)

		return c.Next()
	}
}
