package admin

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"Backend/config"
	"Backend/database"
)

var (
	oauthConfig oauth2.Config
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	oauthState  = "randomStateString"
	initOnce    sync.Once
)

type Admin struct {
	AdminID   string `json:"adminID"`
	EmailID   string `json:"emailID"`
	AdminName string `json:"adminName"`
}

func InitOAuth() {
	initOnce.Do(func() {
		cfg := config.GetConfig()
		if cfg == nil {
			log.Fatal("Failed to load configuration. Ensure .env is correctly set.")
		}

		oauthConfig = oauth2.Config{
			ClientID:     cfg.OAuthConfig.ClientID,
			ClientSecret: cfg.OAuthConfig.ClientSecret,
			RedirectURL:  cfg.OAuthConfig.RedirectURL,
			Scopes:       cfg.OAuthConfig.Scopes,
			Endpoint:     google.Endpoint,
		}

		if err := LoadKeys(); err != nil {
			log.Fatalf("Failed to load encryption keys: %v", err)
		}
	})
}

func LoadKeys() error {
	privBytes, err := os.ReadFile("middleware/encryptionKeys/privateKey.pem")
	if err != nil {
		return fmt.Errorf("error reading private key file: %w", err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return fmt.Errorf("error parsing private key: %w", err)
	}

	pubBytes, err := os.ReadFile("middleware/encryptionKeys/publicKey.pem")
	if err != nil {
		return fmt.Errorf("error reading public key file: %w", err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return fmt.Errorf("error parsing public key: %w", err)
	}

	return nil
}

func HandleGoogleURL(c *fiber.Ctx) error {
	InitOAuth()
	authURL := oauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)
	if authURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate OAuth URL"})
	}

	return c.JSON(fiber.Map{"URL": authURL})
}

func HandleGoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	if state != oauthState {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid OAuth state"})
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Authorization code not found"})
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("[ERROR]: Code exchange failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Code exchange failed"})
	}

	client := oauthConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
	}
	defer userInfoResp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user info"})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get database connection"})
	}

	var admin Admin
	query := `SELECT adminID, emailID, adminName FROM adminData WHERE emailID = $1`
	err = dbConn.QueryRow(query, userInfo.Email).Scan(&admin.AdminID, &admin.EmailID, &admin.AdminName)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin not found"})
		}
		log.Printf("Error querying adminData: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database query error"})
	}

	claimsJWT := jwt.MapClaims{
		"sub":   admin.AdminID,
		"email": admin.EmailID,
		"name":  admin.AdminName,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsJWT)
	tokenString, err := jwtToken.SignedString(privateKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to sign JWT"})
	}

	return c.JSON(fiber.Map{"jwtToken": tokenString, "userRole": "Admin"})
}