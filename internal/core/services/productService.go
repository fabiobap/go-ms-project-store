package services

import (
	"net/http"
	"time"

	dto "github.com/go-ms-project-store/internal/adapters/input/http/dto/product"
	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	"github.com/go-ms-project-store/internal/pkg/pagination"
)

type ProductService interface {
	GetAllProducts(*http.Request) (domain.Products, int64, pagination.DataDBFilter, *errs.AppError)
	CreateProduct(dto.NewProductRequest) (*domain.Product, *errs.AppError)
	FindProductById(int) (*domain.Product, *errs.AppError)
	DeleteProduct(int) (bool, *errs.AppError)
	UpdateProduct(int64, dto.UpdateProductRequest) (*domain.Product, *errs.AppError)
}

type DefaultProductService struct {
	repo domain.ProductRepository
}

func (s DefaultProductService) GetAllProducts(r *http.Request) (domain.Products, int64, pagination.DataDBFilter, *errs.AppError) {
	allowedOrderBy := map[string]bool{
		"id": true, "name": true, "slug": true, "product_id": true, "created_at": true, "updated_at": true,
	}

	filter := pagination.GetBaseFilterParams(r, allowedOrderBy)
	products, totalRows, err := s.repo.FindAll(filter)

	if err != nil {
		logger.Error("Error while finding all products")
		return nil, 0, pagination.DataDBFilter{}, errs.NewUnexpectedError("unexpected database error")
	}

	return products, totalRows, filter, nil
}

func (s DefaultProductService) CreateProduct(req dto.NewProductRequest) (*domain.Product, *errs.AppError) {
	product := domain.NewProduct(req)

	newProduct, err := s.repo.Create(product)
	if err != nil {
		if err.Code != http.StatusUnprocessableEntity {
			return nil, errs.NewUnexpectedError("unexpected database error")
		} else {
			return nil, err
		}
	}

	return newProduct, nil
}

func (s DefaultProductService) UpdateProduct(id int64, req dto.UpdateProductRequest) (*domain.Product, *errs.AppError) {
	product := domain.Product{
		Id:          id,
		Name:        req.Name,
		Slug:        req.Slug,
		CategoryId:  req.CategoryId,
		Amount:      req.Amount,
		Description: req.Description,
		UpdatedAt:   time.Now(),
	}

	newProduct, err := s.repo.Update(product)
	if err != nil {
		if err.Code == http.StatusUnprocessableEntity {
			return nil, err
		} else if err.Code == http.StatusNotFound {
			return nil, err
		}
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return newProduct, nil
}

func (s DefaultProductService) FindProductById(id int) (*domain.Product, *errs.AppError) {
	product, err := s.repo.FindById(id)
	if err != nil {
		if err.Code != http.StatusNotFound {
			return nil, errs.NewUnexpectedError("unexpected database error")
		} else {
			return nil, err
		}
	}

	return product, nil
}

func (s DefaultProductService) DeleteProduct(id int) (bool, *errs.AppError) {
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

func NewProductService(repository domain.ProductRepository) DefaultProductService {
	return DefaultProductService{repo: repository}
}
