package domain

import (
	"time"

	dto "github.com/go-ms-project-store/dto/category"
	"github.com/go-ms-project-store/errs"
	"github.com/go-ms-project-store/helpers"
)

type Category struct {
	Id        int32     `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Categories []Category

type CategoryRepository interface {
	// Save(Category) (*Category, *errs.AppError)
	// Update(Category) (*Category, *errs.AppError)
	// FindById(id int) (*Category, *errs.AppError)
	FindAll() (Categories, *errs.AppError)
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
