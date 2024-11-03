package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-ms-project-store/internal/adapters/input/http/handlers"
	"github.com/go-ms-project-store/internal/adapters/input/http/middlewares"
	"github.com/go-ms-project-store/internal/core/repositories"
	"github.com/go-ms-project-store/internal/core/services"
	"github.com/go-ms-project-store/internal/pkg/db"
)

func Routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middlewares.StoreRoutePattern)

	dbClient := db.GetDBClient()

	categoryRepositoryDB := repositories.NewCategoryRepositoryDB(dbClient)
	ch := handlers.NewCategoryHandlers(services.NewCategoryService(categoryRepositoryDB))

	mux.Route("/api/v1", func(mux chi.Router) {
		mux.Get("/home", handlers.Home)

		mux.Route("/admin", func(mux chi.Router) {
			mux.Get("/categories", ch.GetAllCategories)
		})
	})

	return mux
}
