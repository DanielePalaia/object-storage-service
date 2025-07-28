package main

import (
	"log"
	"os"

	"github.com/DanielePalaia/object-storage-service/api"
	"github.com/DanielePalaia/object-storage-service/persistence"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	storage := persistence.NewInMemoryStorage()
	srv := api.NewServer(storage, port)

	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
