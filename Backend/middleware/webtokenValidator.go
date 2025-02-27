package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
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
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
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
		"sub":   claims["sub"],
		"studentEmail": claims["email"],
		"studentName":  claims["name"],
	}

	body, err := json.Marshal(userData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode user data"})
	}

	c.Request().SetBody(body)
	return c.Next()
}