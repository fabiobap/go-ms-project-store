package services

import (
	"net/http"
	"strconv"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
)

type CategoryService interface {
	GetAllCategories(*http.Request) (domain.Categories, int64, dto.DataDBFilter, *errs.AppError)
	CreateCategory(dto.NewCategoryRequest) (*domain.Category, *errs.AppError)
}

type DefaultCategoryService struct {
	repo domain.CategoryRepository
}

func (s DefaultCategoryService) GetAllCategories(r *http.Request) (domain.Categories, int64, dto.DataDBFilter, *errs.AppError) {
	allowedOrderBy := map[string]bool{
		"id": true, "name": true, "slug": true, "created_at": true, "updated_at": true,
	}

	filter := GetBaseFilterParams(r, allowedOrderBy)
	categories, totalRows, err := s.repo.FindAll(filter)

	if err != nil {
		logger.Error("Error while finding all categories")
		return nil, 0, dto.DataDBFilter{}, errs.NewUnexpectedError("unexpected database error")
	}

	return categories, totalRows, filter, nil
}

func (s DefaultCategoryService) CreateCategory(req dto.NewCategoryRequest) (*domain.Category, *errs.AppError) {
	category := domain.NewCategory(req)

	newCategory, err := s.repo.Create(category)
	if err != nil {
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return newCategory, nil
}

func NewCategoryService(repository domain.CategoryRepository) DefaultCategoryService {
	return DefaultCategoryService{repo: repository}
}

func GetBaseFilterParams(r *http.Request, allowedOrderBy map[string]bool) dto.DataDBFilter {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 15 // default value
	}

	orderBy := r.URL.Query().Get("order_by")
	if orderBy == "" {
		orderBy = "id"
	}
	if !allowedOrderBy[orderBy] {
		orderBy = "id"
	}

	orderDir := r.URL.Query().Get("order_dir")
	if orderDir == "" {
		orderDir = "desc"
	}
	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "desc"
	}

	return dto.DataDBFilter{
		OrderBy:  orderBy,
		OrderDir: orderDir,
		Page:     page,
		PerPage:  perPage,
	}
}
