package main

import (
	"github.com/go-ms-project-store/app"
	"github.com/go-ms-project-store/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	logger.Info("Starting the application")
	app.Start()
}
