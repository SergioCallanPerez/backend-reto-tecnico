package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/model"
)

func Middleware(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return unauthorized(c, "falta el header Authorization")
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return unauthorized(c, "el header Authorization debe ser \"Bearer <token>\"")
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !token.Valid {
			return unauthorized(c, "token inválido o expirado")
		}

		return c.Next()
	}
}

func unauthorized(c *fiber.Ctx, message string) error {
	appErr := model.NewUnauthorizedError(message)
	return c.Status(appErr.HTTPStatus()).JSON(appErr)
}
