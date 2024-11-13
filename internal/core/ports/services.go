package ports

import (
	"net/http"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/pagination"
)

type AuthService interface {
	Login(dto.NewLoginRequest) (*dto.TokenResponse, *errs.AppError)
	Logout(uint64) *errs.AppError
	Me(uint64) (*domain.User, *errs.AppError)
	Register(dto.NewUserRegisterRequest) (*domain.User, *errs.AppError)
}

type CategoryService interface {
	GetAllCategories(*http.Request) (domain.Categories, int64, pagination.DataDBFilter, *errs.AppError)
	CreateCategory(dto.NewCategoryRequest) (*domain.Category, *errs.AppError)
	FindCategoryById(int) (*domain.Category, *errs.AppError)
	DeleteCategory(int) (bool, *errs.AppError)
	UpdateCategory(int64, dto.UpdateCategoryRequest) (*domain.Category, *errs.AppError)
}

type ProductService interface {
	GetAllProducts(*http.Request) (domain.Products, int64, pagination.DataDBFilter, *errs.AppError)
	CreateProduct(dto.NewProductRequest) (*domain.Product, *errs.AppError)
	FindProductById(int) (*domain.Product, *errs.AppError)
	DeleteProduct(int) (bool, *errs.AppError)
	UpdateProduct(int64, dto.UpdateProductRequest) (*domain.Product, *errs.AppError)
}

type UserService interface {
	GetAllUserCustomers(*http.Request) (domain.Users, int64, pagination.DataDBFilter, *errs.AppError)
	GetAllUserAdmins(*http.Request) (domain.Users, int64, pagination.DataDBFilter, *errs.AppError)
	// GetAllUsers(*http.Request) (domain.Users, int64, pagination.DataDBFilter, *errs.AppError)
	FindUserById(string) (*domain.User, *errs.AppError)
	DeleteUser(string) (bool, *errs.AppError)
}
