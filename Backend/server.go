package main

import (
	"Backend/config"
	"Backend/database"
	// "Backend/routes"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func main() {

	fmt.Println("[LOG]: Starting Virtual Classroom Backend Server..");

	cfg := config.GetConfig();
	portNum := ":" + cfg.ServerPort;

	config.LoadEnv();
	database.ConnectDB();

	var wg sync.WaitGroup
	wg.Add(1);
	go database.RunMigrations(&wg);
	wg.Wait();

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