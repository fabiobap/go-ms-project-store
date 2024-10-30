package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-ms-project-store/app/handlers"
)

func Start() {
	mux := chi.NewRouter()

	// dbClient := getDBClient()
	mux.Get("/home", handlers.Home)

	log.Fatal(http.ListenAndServe("localhost:8686", mux))
}
