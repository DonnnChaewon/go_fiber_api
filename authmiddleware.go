package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c * fiber.Ctx) error {
		cookieToken := c.Cookies("jwt")
		var tokenString string

		if cookieToken != "" {
			log.Println("Token from cookies, using it...")
			tokenString = cookieToken
		} else {
			log.Println("Empty token from cookies, trying to get it from the auth header...")
			authHeader := c.Get("Authorization")

			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			tokenParts := strings.Split(authHeader, " ")

			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			tokenString = tokenParts[1]
		}

		secret := []byte("super-secret-key")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		userId := token.Claims.(jwt.MapClaims)["userId"]

		if err := db.Model(&User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		c.Locals("userId", userId)

		return c.Next()
	}
}