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
	productRepositoryDB := repositories.NewProductRepositoryDB(dbClient)
	ch := handlers.NewCategoryHandlers(services.NewCategoryService(categoryRepositoryDB))
	ph := handlers.NewProductHandlers(services.NewProductService(productRepositoryDB))

	mux.Route("/api/v1", func(mux chi.Router) {
		mux.Get("/home", handlers.Home)

		mux.Route("/admin", func(mux chi.Router) {
			mux.Route("/categories", func(mux chi.Router) {
				mux.Get("/", ch.GetAllCategories)
				mux.Get("/{id}", ch.GetCategory)
				mux.Post("/", ch.CreateCategory)
				mux.Put("/{id}", ch.UpdateCategory)
				mux.Delete("/{id}", ch.DeleteCategory)
			})
			mux.Route("/products", func(mux chi.Router) {
				mux.Get("/", ph.GetAllProducts)
				mux.Get("/{id}", ph.GetProduct)
				mux.Post("/", ph.CreateProduct)
				mux.Put("/{id}", ph.UpdateProduct)
				mux.Delete("/{id}", ph.DeleteProduct)
			})
		})
	})

	return mux
}
