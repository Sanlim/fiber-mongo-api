package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

var (
	jwtSecret = os.Getenv("JWT_SECRET")
)

//https://github.com/gofiber/jwt
func AuthorizationRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		// Filter:         nil,
		SuccessHandler: AuthSuccess,
		ErrorHandler:   AuthError,
		SigningKey:     []byte(jwtSecret),
		// SigningKeys:   nil,
		SigningMethod: "HS256",
		// ContextKey:    nil,
		// Claims:        nil,
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
	})
}

func AuthError(c *fiber.Ctx, e error) error {
	c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
		"msg":   e.Error(),
	})
	return nil
}

func AuthSuccess(c *fiber.Ctx) error {
	c.Next()
	return nil
}
