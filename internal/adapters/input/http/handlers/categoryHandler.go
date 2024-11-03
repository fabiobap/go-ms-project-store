package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/services"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type CategoryHandlers struct {
	Service services.CategoryService
}

func NewCategoryHandlers(service services.CategoryService) *CategoryHandlers {
	return &CategoryHandlers{
		Service: service,
	}
}

func (ch *CategoryHandlers) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, totalRows, filter, err := ch.Service.GetAllCategories(r)

	baseURL := helpers.GetFullRouteUrl(r)

	paginatedResponse := dto.NewPaginatedResponse(categories.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, paginatedResponse)
	}
}

func (ch *CategoryHandlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryRequest dto.NewCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&categoryRequest)
	if err != nil {
		helpers.WriteResponse(w, http.StatusBadRequest, err.Error())
	}

	if err := dto.ValidateCategory(&categoryRequest); err != nil {
		helpers.WriteResponse(w, http.StatusUnprocessableEntity, err)
	}

	category, errCat := ch.Service.CreateCategory(categoryRequest)
	if errCat != nil {
		helpers.WriteResponse(w, errCat.Code, errCat.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusCreated, category.ToCategoryDTO())
	}
}
