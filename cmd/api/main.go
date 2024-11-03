package main

import (
	"log"
	"net/http"

	"github.com/go-ms-project-store/internal/adapters/input/http/routes"
	"github.com/go-ms-project-store/internal/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	logger.Info("Starting the application")
	mux := routes.Routes()

	log.Fatal(http.ListenAndServe("localhost:8686", mux))

}
