package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type Category struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Categories []Category

type CategoryRepository interface {
	Create(Category) (*Category, *errs.AppError)
	Delete(id int) *errs.AppError
	FindAll(filter dto.DataDBFilter) (Categories, int64, *errs.AppError)
	FindById(id int) (*Category, *errs.AppError)
	// Update(Category) (*Category, *errs.AppError)
}

func NewCategory(req dto.NewCategoryRequest) Category {
	return Category{
		Name:      req.Name,
		Slug:      req.Slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c Category) ToCategoryDTO() dto.CategoryResponse {
	return dto.CategoryResponse{
		Id:        c.Id,
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: helpers.DatetimeToString(c.CreatedAt),
		UpdatedAt: helpers.DatetimeToString(c.UpdatedAt),
	}
}

func (c Categories) ToDTO() []dto.CategoryResponse {
	dtos := make([]dto.CategoryResponse, len(c))
	for i, category := range c {
		dtos[i] = category.ToCategoryDTO()
	}
	return dtos
}
