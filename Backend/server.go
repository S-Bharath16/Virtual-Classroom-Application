package main

import (
	"Backend/config"
	"Backend/database"
	"Backend/routes"
	"Backend/utilities/RSA"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}

func main() {
	fmt.Println("[LOG]: Starting Virtual Classroom Backend Server..")

	cfg := config.GetConfig()
	portNum := ":" + cfg.ServerPort

	config.LoadEnv()
	database.ConnectDB()

	privateKeyPath := "middleware/encryptionKeys/privateKey.pem"
	publicKeyPath := "middleware/encryptionKeys/publicKey.pem"

	privateKeyExists := fileExists(privateKeyPath)
	publicKeyExists := fileExists(publicKeyPath)

	if !privateKeyExists || !publicKeyExists {
		fmt.Println("[LOG]: Key pair missing! Generating new SSH keys...")

		err := RSA.GenerateRSAKeys(privateKeyPath, publicKeyPath)
		if err != nil {
			log.Fatalf("[ERROR]: Error generating keys: %v", err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go database.RunMigrations(&wg)
	wg.Wait()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001", // ✅ Specify frontend URL
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Authorization,Content-Type",
		AllowCredentials: true,
	}))

	routes.RegisterStudent(app)
	routes.RegisterFacultyRoutes(app)
	routes.RegisterStudentRoutes(app)

	fmt.Println("[LOG]: Server Started on Port:", portNum)
	log.Fatal(app.Listen(portNum)) // ✅ Use dynamic port from config

	fmt.Println("[LOG]: To close connection CTRL+C :-)")
}