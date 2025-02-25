package config

import (
	"log"
	"os"
	"strings"
	"github.com/joho/godotenv"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Tenant       string
}

type Config struct {
	DBUrl           string
	ServerPort      string
	Environment     string
	EmailID         string
	MailerPassword  string
	MailHost        string
	MailHostAddress string
	OAuthConfig     OAuthConfig
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func GetConfig() Config {
	LoadEnv()

	return Config{
		DBUrl:           os.Getenv("DB_URL"),
		ServerPort:      os.Getenv("SERVER_PORT"),
		Environment:     os.Getenv("ENVIRONMENT"),
		EmailID:         os.Getenv("EMAIL_ID"),
		MailerPassword:  os.Getenv("MAILER_PWD"),
		MailHost:        os.Getenv("MAIL_HOST"),
		MailHostAddress: os.Getenv("MAIL_HOST_ADDRESS"),
		OAuthConfig: OAuthConfig{
			ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
			Scopes:       splitAndTrim(os.Getenv("OAUTH_SCOPES"), ","),
			Tenant:       os.Getenv("OAUTH_TENANT"),
		},
	}
}