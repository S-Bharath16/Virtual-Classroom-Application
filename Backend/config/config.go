package config

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/joho/godotenv"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName 		string
	Config          aws.Config
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
	AWSConfig       AWSConfig
}

var (
	configInstance Config
	once           sync.Once
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Could not load .env file, falling back to system environment variables", err)
	}
}

func loadAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"", 
		)),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	return cfg
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
			AWSConfig: AWSConfig{
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				Region:          os.Getenv("AWS_REGION"),
				BucketName: 	 os.Getenv("AWS_BUCKET_NAME"),
				Config:          loadAWSConfig(),
			},
		}
	})
	return &configInstance
}