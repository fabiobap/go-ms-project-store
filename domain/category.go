package domain

import (
	"time"

	"github.com/go-ms-project-store/errs"
)

type Category struct {
	Id        int32     `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CategoryRepository interface {
	// Save(Category) (*Category, *errs.AppError)
	// Update(Category) (*Category, *errs.AppError)
	// FindById(id int) (*Category, *errs.AppError)
	FindAll() ([]Category, *errs.AppError)
}
