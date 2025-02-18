package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl       string
	ServerPort  string
	Environment string
	EmailID 	string
	MailerPassword string
	MailHost 	string
	MailHostAddress string
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetConfig() Config {
	LoadEnv()

	return Config{
		DBUrl:       os.Getenv("DB_URL"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
		EmailID: os.Getenv("EMAIL_ID"),
		MailerPassword: os.Getenv("MAILER_PWD"),
		MailHost: os.Getenv("MAIL_HOST"),
		MailHostAddress: os.Getenv("MAIL_HOST_ADDRESS"),
	}
}