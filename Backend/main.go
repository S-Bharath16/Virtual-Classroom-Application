package main

import (
	"fmt"
	"net/http"
	"log"
)


func main() {
	const portNum string = ":8080"

	log.Println("Starting our simple http server.")

    log.Println("Started on port", portNum)
    fmt.Println("To close connection CTRL+C :-)")

    err := http.ListenAndServe(portNum, nil)
    if err != nil {
        log.Fatal(err)
    }
}