package service

import (
	domain "github.com/go-ms-project-store/domain/category"
	"github.com/go-ms-project-store/errs"
	"github.com/go-ms-project-store/logger"
)

type CategoryService interface {
	GetAllCategories() (domain.Categories, *errs.AppError)
	// GetCategory(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCategoryService struct {
	repo domain.CategoryRepository
}

func (s DefaultCategoryService) GetAllCategories() (domain.Categories, *errs.AppError) {
	categories, err := s.repo.FindAll()
	if err != nil {
		logger.Error("Error while finding all categories")
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return categories, nil
}

func NewCategoryService(repository domain.CategoryRepository) DefaultCategoryService {
	return DefaultCategoryService{repo: repository}
}
