package Auth

import (
	"os"
	"fmt"
	"log"
	"time"
	"context"
	"crypto/rsa"
	"database/sql"

	"Backend/config"
	"Backend/models"
	"Backend/database"
	
	"golang.org/x/oauth2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2/microsoft"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	oauthState = "randomStateString"
	oauthConf  *oauth2.Config
)

func LoadKeys() error {
	privBytes, err := os.ReadFile("middleware/encryptionKeys/privateKey.pem")
	if err != nil {
		return fmt.Errorf("error reading private key file: %w", err)
	}
	pKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return fmt.Errorf("error parsing private key: %w", err)
	}
	privateKey = pKey

	pubBytes, err := os.ReadFile("middleware/encryptionKeys/publicKey.pem")
	if err != nil {
		return fmt.Errorf("error reading public key file: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return fmt.Errorf("error parsing public key: %w", err)
	}
	publicKey = pubKey

	return nil
}

func init() {

	if err := LoadKeys(); err != nil {
		log.Fatalf("Failed to load RSA keys: %v", err)
	}

	cfg := config.GetConfig()

	oauthConf = &oauth2.Config{
		ClientID:     cfg.OAuthConfig.ClientID,
		ClientSecret: cfg.OAuthConfig.ClientSecret,
		RedirectURL:  cfg.OAuthConfig.RedirectURL,
		Scopes:       cfg.OAuthConfig.Scopes,
		Endpoint:     microsoft.AzureADEndpoint(cfg.OAuthConfig.Tenant),
	}
}

func HandleMicrosoftURL(c *fiber.Ctx) error {
	authURL := oauthConf.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)
	return c.JSON(fiber.Map{"URL": authURL})
}

func HandleMicrosoftCallback(c *fiber.Ctx) error {

	state := c.Query("state")
	if state != oauthState {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid OAuth state"})
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Authorization code not found"})
	}

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("[ERROR]: Code exchange failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Code exchange failed"})
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No id_token found in token response"})
	}

	parsedToken, _, err := new(jwt.Parser).ParseUnverified(rawIDToken, jwt.MapClaims{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse id_token"})
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to extract claims from id_token"})
	}

	email, ok := claims["preferred_username"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Email not found in token claims"})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("[ERROR]: Error getting DB connection: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get database connection"})
	}

	var student models.Student
	rawQuery := `
        SELECT studentID, rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID 
        FROM studentData 
        WHERE emailID = $1
    `
	err = dbConn.QueryRow(rawQuery, email).Scan(
		&student.StudentID,
		&student.RollNumber,
		&student.EmailID,
		&student.StudentName,
		&student.StartYear,
		&student.EndYear,
		&student.DeptID,
		&student.SectionID,
		&student.SemesterID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Student not found"})
		}
		log.Printf("[ERROR]: Error querying studentData: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database query error"})
	}

	claimsJWT := jwt.MapClaims{
		"sub":   student.StudentID,
		"email": student.EmailID,
		"name":  student.StudentName,
		"role":  "Student",
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	}
	
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsJWT)
	tokenString, err := jwtToken.SignedString(privateKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to sign JWT"})
	}

	return c.JSON(fiber.Map{"JWT Token": tokenString})
}