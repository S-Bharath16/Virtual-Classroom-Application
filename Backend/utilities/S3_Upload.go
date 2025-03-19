package utilities

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appConfig "Backend/config"
)

var (
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
)

func NewS3Clients() {
	cfg := appConfig.GetConfig().AWSConfig.Config

	S3Client = s3.NewFromConfig(cfg)
	PresignClient = s3.NewPresignClient(S3Client)
}

func PutPresignURL(key string) string {
	cfg := appConfig.GetConfig().AWSConfig

	presignedUrl, err := PresignClient.PresignPutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(cfg.BucketName),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(10 * time.Minute),
	)
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		return ""
	}
	return presignedUrl.URL
}

func UploadFile(url, filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(file))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	log.Println("File uploaded successfully")
	return nil
}

func DeleteS3Object(key string) bool {
	cfg := appConfig.GetConfig().AWSConfig

	_, err := S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(cfg.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("Failed to delete object: %v", err)
		return false
	}
	log.Println("File deleted successfully")
	return true
}

// ✅ Example Usage (Optional)
func ExampleUsage() {
	// Step 1: Load AWS Config from config package
	appConfig.GetConfig() // This ensures AWS config is initialized

	// Step 2: Initialize S3 Client
	NewS3Clients()

	// Step 3: Generate Presigned URL
	url := PutPresignURL("example.txt")
	if url == "" {
		log.Println("Failed to generate presigned URL")
		return
	}

	log.Println("Presigned URL:", url)

	// Step 4: Upload File
	err := UploadFile(url, "./example.txt")
	if err != nil {
		log.Printf("Upload failed: %v", err)
	}

	// Step 5: Delete File (Optional)
	success := DeleteS3Object("example.txt")
	if success {
		log.Println("File deleted successfully")
	} else {
		log.Println("Failed to delete file")
	}
}