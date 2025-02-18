package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	// "sync"
	"Backend/config"
	// "Backend/database"
	// "Backend/routes"
	"Backend/utilities/RSA"
	"github.com/gofiber/fiber/v2"
)

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}

func main() {

	fmt.Println("[LOG]: Starting Virtual Classroom Backend Server..");

	cfg := config.GetConfig();
	portNum := ":" + cfg.ServerPort;

	config.LoadEnv();
	// database.ConnectDB();

	privateKeyPath := "middleware/encryptionKeys/privateKey.pem";
	publicKeyPath := "middleware/encryptionKeys/publicKey.pem";

	privateKeyExists := fileExists(privateKeyPath);
	publicKeyExists := fileExists(publicKeyPath);

	if !privateKeyExists || !publicKeyExists {
		fmt.Println("[LOG]: Key pair missing! Generating new SSH keys...");

		err := RSA.GenerateRSAKeys(privateKeyPath, publicKeyPath);
		if err != nil {
			log.Fatalf("[ERROR]: Error generating keys: %v", err);
		}
	}

	// var wg sync.WaitGroup
	// wg.Add(1);
	// go database.RunMigrations(&wg);
	// wg.Wait();

	app := fiber.New();

	// routes.RegisterItemRoutes(app);

	fmt.Println("[LOG]: Server Started on Port: ", portNum)
	log.Fatal(app.Listen(":8080"))
	fmt.Println("[LOG]: To close connection CTRL+C :-)")

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}