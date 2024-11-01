package service

import (
	domain "github.com/go-ms-project-store/domain/category"
	"github.com/go-ms-project-store/errs"
)

type CategoryService interface {
	GetAllCategories() ([]domain.Category, *errs.AppError)
	// GetCategory(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCategoryService struct {
	repo domain.CategoryRepository
}

func (s DefaultCategoryService) GetAllCategories() ([]domain.Category, *errs.AppError) {
	return s.repo.FindAll()
}

func NewCategoryService(repository domain.CategoryRepository) DefaultCategoryService {
	return DefaultCategoryService{repo: repository}
}
