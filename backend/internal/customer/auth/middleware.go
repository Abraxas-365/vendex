package auth

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

// CustomerMiddleware provides Fiber middleware for customer JWT authentication.
type CustomerMiddleware struct {
	jwtSecret []byte
}

// NewCustomerMiddleware creates a new customer authentication middleware.
func NewCustomerMiddleware(jwtSecret string) *CustomerMiddleware {
	return &CustomerMiddleware{
		jwtSecret: []byte(jwtSecret),
	}
}

// Authenticate returns a Fiber handler that validates the customer Bearer JWT
// and stores the parsed context in c.Locals("customer_auth").
func (m *CustomerMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		var tokenStr string

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" && parts[1] != "" {
				tokenStr = parts[1]
			}
		}

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   ErrTokenRequired().Code,
				"message": ErrTokenRequired().Message,
			})
		}

		claims, err := m.parseToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   ErrTokenInvalid().Code,
				"message": ErrTokenInvalid().Message,
			})
		}

		c.Locals("customer_auth", &CustomerAuthContext{
			CustomerID: claims.CustomerID,
			TenantID:   claims.TenantID,
			Email:      claims.Email,
			Name:       claims.Name,
		})

		return c.Next()
	}
}

// GetCustomerAuthContext extracts the customer auth context from Fiber locals.
func GetCustomerAuthContext(c *fiber.Ctx) (*CustomerAuthContext, bool) {
	auth, ok := c.Locals("customer_auth").(*CustomerAuthContext)
	return auth, ok && auth != nil
}

func (m *CustomerMiddleware) parseToken(tokenStr string) (*CustomerClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomerClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid()
	}
	claims, ok := token.Claims.(*CustomerClaims)
	if !ok {
		return nil, ErrTokenInvalid()
	}
	return claims, nil
}
