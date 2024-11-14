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

	dbClient := db.GetDBClient()

	authRepositoryDB := repositories.NewAuthRepositoryDB(dbClient)

	mux.Use(middlewares.StoreRoutePattern)
	authMiddleware := middlewares.NewAuthMiddleware(authRepositoryDB)

	categoryRepositoryDB := repositories.NewCategoryRepositoryDB(dbClient)
	productRepositoryDB := repositories.NewProductRepositoryDB(dbClient)
	userRepositoryDB := repositories.NewUserRepositoryDB(dbClient)

	ah := handlers.NewAuthHandlers(services.NewAuthService(authRepositoryDB))
	ch := handlers.NewCategoryHandlers(services.NewCategoryService(categoryRepositoryDB))
	ph := handlers.NewProductHandlers(services.NewProductService(productRepositoryDB))
	uh := handlers.NewUserHandlers(services.NewUserService(userRepositoryDB))

	mux.Route("/api/v1", func(mux chi.Router) {
		mux.Get("/home", handlers.Home)

		mux.Route("/auth", func(mux chi.Router) {
			mux.Post("/login", ah.Login)
			mux.Post("/register", ah.Register)
			mux.Group(func(mux chi.Router) {
				mux.Use(authMiddleware.Auth)
				mux.Post("/refresh-token", ah.Refresh)
				mux.Post("/logout", ah.Logout)
				mux.Get("/me", ah.Me)
			})
		})
		mux.Route("/admin", func(mux chi.Router) {
			mux.Use(authMiddleware.Auth)
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
			mux.Route("/users", func(mux chi.Router) {
				mux.Get("/user-admins", uh.GetAllUserAdmins)
				mux.Get("/user-customers", uh.GetAllUserCustomers)
				mux.Get("/{id}", uh.GetUser)
				mux.Delete("/{id}", uh.DeleteUser)
			})
		})
	})

	return mux
}
