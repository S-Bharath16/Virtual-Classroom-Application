package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
)

var publicKey *rsa.PublicKey

func init() {
	keyData, err := os.ReadFile("middleware/encryptionKeys/publicKey.pem")
	if err != nil {
		panic("failed to load public key: " + err.Error())
	}

	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		panic("failed to parse public key: " + err.Error())
	}

	publicKey = parsedKey
}

func WebTokenValidator(c *fiber.Ctx) error {
	if !strings.HasPrefix(c.Path(), "/api/student/") && !strings.HasPrefix(c.Path(), "/api/admin") {
		return c.Next()
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	tokenString := parts[1]
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is empty"})
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	userData := fiber.Map{
		"sub":      claims["sub"],
		"emailID":  claims["email"],
		"userName": claims["name"],
	}

	body, err := json.Marshal(userData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode user data"})
	}

	c.Request().SetBody(body)
	return c.Next()
}