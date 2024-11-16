package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/helpers"
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

func NewProduct(req dto.NewProductRequest) Product {
	return Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Amount:      req.Amount,
		Image:       "https://placehold.co/600x400",
		CategoryId:  req.CategoryId,
		UUID:        uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (p Product) ToProductDTO() dto.ProductResponse {
	amount := float64(p.Amount) / 100
	return dto.ProductResponse{
		Id:          p.Id,
		UUID:        p.UUID,
		Name:        p.Name,
		Description: p.Description,
		Amount:      helpers.NumberFormat(amount, 2, ".", ","),
		Image:       p.Image,
		Slug:        p.Slug,
		CategoryId:  p.CategoryId,
		CreatedAt:   helpers.DatetimeToString(p.CreatedAt),
		UpdatedAt:   helpers.DatetimeToString(p.UpdatedAt),
		Category: dto.CategoryResponse{
			Id:   p.Category.Id,
			Name: p.Category.Name,
			Slug: p.Category.Slug,
			CreatedAt: helpers.DatetimeToString(
				p.Category.CreatedAt,
			),
			UpdatedAt: helpers.DatetimeToString(
				p.Category.UpdatedAt,
			),
		},
	}
}

func (p Product) ToPublicProductDTO() dto.ProductPublicResponse {
	amount := float64(p.Amount) / 100
	return dto.ProductPublicResponse{
		ID:          p.UUID,
		Name:        p.Name,
		Description: p.Description,
		Amount:      helpers.NumberFormat(amount, 2, ".", ","),
		Image:       p.Image,
		Slug:        p.Slug,
		CreatedAt:   helpers.DatetimeToString(p.CreatedAt),
		Category: dto.CategoryPublicResponse{
			Name: p.Category.Name,
			Slug: p.Category.Slug,
			CreatedAt: helpers.DatetimeToString(
				p.Category.CreatedAt,
			),
		},
	}
}

func (c Products) ToDTO() []dto.ProductResponse {
	dtos := make([]dto.ProductResponse, len(c))
	for i, product := range c {
		dtos[i] = product.ToProductDTO()
	}
	return dtos
}

func (c Products) ToPublicDTO() []dto.ProductPublicResponse {
	dtos := make([]dto.ProductPublicResponse, len(c))
	for i, product := range c {
		dtos[i] = product.ToPublicProductDTO()
	}
	return dtos
}
