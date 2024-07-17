package main

import (
	"fmt"
	"future-app/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiServer := server.NewAPIServer(fmt.Sprintf(":%s", port))
	apiServer.Run()
}
