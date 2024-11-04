package services

import (
	"net/http"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	"github.com/go-ms-project-store/internal/pkg/pagination"
)

type CategoryService interface {
	GetAllCategories(*http.Request) (domain.Categories, int64, pagination.DataDBFilter, *errs.AppError)
	CreateCategory(dto.NewCategoryRequest) (*domain.Category, *errs.AppError)
	FindCategoryById(int) (*domain.Category, *errs.AppError)
	DeleteCategory(int) (bool, *errs.AppError)
}

type DefaultCategoryService struct {
	repo domain.CategoryRepository
}

func (s DefaultCategoryService) GetAllCategories(r *http.Request) (domain.Categories, int64, pagination.DataDBFilter, *errs.AppError) {
	allowedOrderBy := map[string]bool{
		"id": true, "name": true, "slug": true, "created_at": true, "updated_at": true,
	}

	filter := pagination.GetBaseFilterParams(r, allowedOrderBy)
	categories, totalRows, err := s.repo.FindAll(filter)

	if err != nil {
		logger.Error("Error while finding all categories")
		return nil, 0, pagination.DataDBFilter{}, errs.NewUnexpectedError("unexpected database error")
	}

	return categories, totalRows, filter, nil
}

func (s DefaultCategoryService) CreateCategory(req dto.NewCategoryRequest) (*domain.Category, *errs.AppError) {
	category := domain.NewCategory(req)

	newCategory, err := s.repo.Create(category)
	if err != nil {
		if err.Code != http.StatusUnprocessableEntity {
			return nil, errs.NewUnexpectedError("unexpected database error")
		} else {
			return nil, err
		}
	}

	return newCategory, nil
}

func (s DefaultCategoryService) FindCategoryById(id int) (*domain.Category, *errs.AppError) {
	category, err := s.repo.FindById(id)
	if err != nil {
		if err.Code != http.StatusNotFound {
			return nil, errs.NewUnexpectedError("unexpected database error")
		} else {
			return nil, err
		}
	}

	return category, nil
}

func (s DefaultCategoryService) DeleteCategory(id int) (bool, *errs.AppError) {
	err := s.repo.Delete(id)
	if err != nil {
		if err.Code != http.StatusNotFound {
			return false, errs.NewUnexpectedError("unexpected database error")
		} else {
			return false, err
		}
	}

	return true, nil
}

func NewCategoryService(repository domain.CategoryRepository) DefaultCategoryService {
	return DefaultCategoryService{repo: repository}
}
