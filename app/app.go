package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-ms-project-store/app/handlers"
	"github.com/go-ms-project-store/domain"
	"github.com/go-ms-project-store/service"
)

func Start() {
	mux := chi.NewRouter()

	dbClient := getDBClient()
	categoryRepositoryDB := domain.NewCategoryRepositoryDB(dbClient)

	ch := CategoryHandlers{service: service.NewCategoryService(categoryRepositoryDB)}

	mux.Route("/api/v1", func(mux chi.Router) {
		mux.Get("/home", handlers.Home)

		mux.Route("/admin", func(mux chi.Router) {
			mux.Get("/categories", ch.getAllCategories)
		})
	})

	log.Fatal(http.ListenAndServe("localhost:8686", mux))
}
