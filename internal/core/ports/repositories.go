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
	RoleRepo() RoleRepository
	Register(domain.UserRegister) (*domain.User, *errs.AppError)
	RevokeAccessToken(uint64) *errs.AppError
}

type CategoryRepository interface {
	Create(domain.Category) (*domain.Category, *errs.AppError)
	Delete(int) *errs.AppError
	FindAll(pagination.DataDBFilter) (domain.Categories, int64, *errs.AppError)
	FindById(int) (*domain.Category, *errs.AppError)
	Update(domain.Category) (*domain.Category, *errs.AppError)
}

type OrderRepository interface {
	Create(domain.Order) (*domain.Order, *errs.AppError)
	FindById(uint64) (*domain.Order, *errs.AppError)
	ProductRepo() ProductRepository
	OrderItemRepo() OrderItemRepository
}

type OrderItemRepository interface {
	Create(domain.OrderItem) (*domain.OrderItem, *errs.AppError)
}

type ProductRepository interface {
	Create(domain.Product) (*domain.Product, *errs.AppError)
	Delete(int) *errs.AppError
	FindAll(pagination.DataDBFilter) (domain.Products, int64, *errs.AppError)
	FindById(int) (*domain.Product, *errs.AppError)
	FindBySlug(string) (*domain.Product, *errs.AppError)
	Update(domain.Product) (*domain.Product, *errs.AppError)
	WhereIn([]string) ([]domain.Product, *errs.AppError)
}

type RoleRepository interface {
	FindByName(string) (*domain.Role, *errs.AppError)
}

type UserRepository interface {
	Delete(string) *errs.AppError
	FindAll(pagination.DataDBFilter, string) (domain.Users, int64, *errs.AppError)
	FindAllAdmins(pagination.DataDBFilter) (domain.Users, int64, *errs.AppError)
	FindAllCustomers(pagination.DataDBFilter) (domain.Users, int64, *errs.AppError)
	FindById(uint64) (*domain.User, *errs.AppError)
	FindByUuid(string) (*domain.User, *errs.AppError)
}
