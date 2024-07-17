package main

import (
	"fmt"
	"future-app/server"
	"future-app/store"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbStore, err := store.NewStore()
	if err != nil {
		log.Fatalf("Error creating store: %v", err)
	}
	defer dbStore.Close()
	if err := dbStore.Init(); err != nil {
		log.Fatalf("Error initializing store: %v", err)
	}

	apiServer := server.NewAPIServer(fmt.Sprintf(":%s", port), dbStore)
	apiServer.Run()
}
