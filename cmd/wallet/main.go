package main

import (
	"log"

	"github.com/pyuldashev912/alif-task-go/internal/server"
)

func main() {

	// Since this is a test task, I decided to leave the configurations here.
	// In the real problem I would read from the environment variables.
	config := server.Config{
		BindAddr:    ":8080",
		DatabaseURL: "host=localhost port=5432 user=postgres dbname=postgres password=alif sslmode=disable",
		CacheAddr:   "localhost:6379",
	}

	if err := server.Start(&config); err != nil {
		log.Fatal(err)
	}
}
