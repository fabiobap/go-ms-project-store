package handlers

import (
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
	baseURL := helpers.GetBaseURL(r) + "/api/v1/admin/categories"

	paginatedResponse := dto.NewPaginatedResponse(categories.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, paginatedResponse)
	}
}
