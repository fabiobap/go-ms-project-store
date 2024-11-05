package domain

import (
	"time"

	dto "github.com/go-ms-project-store/internal/adapters/input/http/dto/product"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/go-ms-project-store/internal/pkg/pagination"
	"github.com/google/uuid"
)

type Product struct {
	Id          int64     `db:"id"`
	UUID        uuid.UUID `db:"uuid"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Amount      int32     `db:"amount"`
	Image       string    `db:"image"`
	Slug        string    `db:"slug"`
	CategoryId  int64     `db:"category_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Category    Category
}

type Products []Product

type ProductRepository interface {
	Create(Product) (*Product, *errs.AppError)
	Delete(id int) *errs.AppError
	FindAll(filter pagination.DataDBFilter) (Products, int64, *errs.AppError)
	FindById(id int) (*Product, *errs.AppError)
	Update(Product) (*Product, *errs.AppError)
}

func NewProduct(req dto.NewProductRequest) Product {
	return Product{
		Name:      req.Name,
		Slug:      req.Slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c Product) ToProductDTO() dto.ProductResponse {
	return dto.ProductResponse{
		Id:          c.Id,
		UUID:        c.UUID,
		Name:        c.Name,
		Description: c.Description,
		Amount:      c.Amount,
		Image:       c.Image,
		Slug:        c.Slug,
		CategoryId:  c.CategoryId,
		CreatedAt:   helpers.DatetimeToString(c.CreatedAt),
		UpdatedAt:   helpers.DatetimeToString(c.UpdatedAt),
	}
}

func (c Products) ToDTO() []dto.ProductResponse {
	dtos := make([]dto.ProductResponse, len(c))
	for i, product := range c {
		dtos[i] = product.ToProductDTO()
	}
	return dtos
}
