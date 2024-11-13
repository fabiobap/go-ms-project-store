package ports

import (
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/pagination"
)

type AuthRepository interface {
	CreateAccessToken(domain.Token) (*domain.Token, *errs.AppError)
	CreateRefreshToken(domain.Token) (*domain.Token, *errs.AppError)
	ValidateToken(string) (uint64, *errs.AppError)
	Login(domain.AuthUser) (*domain.User, *errs.AppError)
	Logout(uint64) *errs.AppError
	UserRepo() UserRepository
	// Register(User) (*User, *errs.AppError)
}

type CategoryRepository interface {
	Create(domain.Category) (*domain.Category, *errs.AppError)
	Delete(id int) *errs.AppError
	FindAll(filter pagination.DataDBFilter) (domain.Categories, int64, *errs.AppError)
	FindById(id int) (*domain.Category, *errs.AppError)
	Update(domain.Category) (*domain.Category, *errs.AppError)
}

type ProductRepository interface {
	Create(domain.Product) (*domain.Product, *errs.AppError)
	Delete(id int) *errs.AppError
	FindAll(filter pagination.DataDBFilter) (domain.Products, int64, *errs.AppError)
	FindById(id int) (*domain.Product, *errs.AppError)
	Update(domain.Product) (*domain.Product, *errs.AppError)
}

type UserRepository interface {
	Delete(id string) *errs.AppError
	FindAll(filter pagination.DataDBFilter, roleName string) (domain.Users, int64, *errs.AppError)
	FindAllAdmins(filter pagination.DataDBFilter) (domain.Users, int64, *errs.AppError)
	FindAllCustomers(filter pagination.DataDBFilter) (domain.Users, int64, *errs.AppError)
	FindById(id uint64) (*domain.User, *errs.AppError)
	FindByUuid(id string) (*domain.User, *errs.AppError)
}
