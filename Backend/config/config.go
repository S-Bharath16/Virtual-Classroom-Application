package config

import (
	"log"
	"os"
	"sync"
	"github.com/joho/godotenv"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
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

var (
	configInstance Config
	once           sync.Once
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Could not load .env file, falling back to system environment variables")
	}
}

func GetConfig() *Config {
	once.Do(func() {
		LoadEnv()

		configInstance = Config{
			DBUrl:           os.Getenv("DB_URL"),
			ServerPort:      os.Getenv("SERVER_PORT"),
			Environment:     os.Getenv("ENVIRONMENT"),
			EmailID:         os.Getenv("EMAIL_ID"),
			MailerPassword:  os.Getenv("MAILER_PWD"),
			MailHost:        os.Getenv("MAIL_HOST"),
			MailHostAddress: os.Getenv("MAIL_HOST_ADDRESS"),
			OAuthConfig: OAuthConfig{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
				Scopes:       []string{"openid", "email", "profile"},
			},
		}
	})
	return &configInstance
}