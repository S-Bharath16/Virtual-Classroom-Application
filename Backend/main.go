package main

import (
	"Backend/config"
	"Backend/database"
	"Backend/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	const portNum string = ":8080"

	config.LoadEnv()
	database.ConnectDB()

	app := fiber.New()

	routes.RegisterItemRoutes(app)

	log.Println("Starting our simple http server.")

	log.Println("Started on port", portNum)
	log.Fatal(app.Listen(":8080"))
	fmt.Println("To close connection CTRL+C :-)")

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
