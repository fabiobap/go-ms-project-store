package app

import (
	"net/http"

	dto "github.com/go-ms-project-store/dto/category"
	"github.com/go-ms-project-store/helpers"
	"github.com/go-ms-project-store/service"
)

type CategoryHandlers struct {
	service service.CategoryService
}

func (ch *CategoryHandlers) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, totalRows, filter, err := ch.service.GetAllCategories(r)
	baseURL := helpers.GetBaseURL(r) + "/api/v1/admin/categories"

	paginatedResponse := dto.NewPaginatedResponse(categories.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, paginatedResponse)
	}
}
