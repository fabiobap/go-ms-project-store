package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/services"
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/go-ms-project-store/internal/pkg/pagination"
)

type CategoryHandlers struct {
	Service services.CategoryService
}

func (ch *CategoryHandlers) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	_, err := ch.Service.DeleteCategory(id)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusNoContent, "")
	}
}

func (ch *CategoryHandlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryRequest dto.NewCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&categoryRequest)
	if err != nil {
		helpers.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateCategory(&categoryRequest); err != nil {
		helpers.WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	category, errCat := ch.Service.CreateCategory(categoryRequest)
	if errCat != nil {
		helpers.WriteResponse(w, errCat.Code, errCat)
	} else {
		helpers.WriteResponse(w, http.StatusCreated, category.ToCategoryDTO())
	}
}

func (ch *CategoryHandlers) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, totalRows, filter, err := ch.Service.GetAllCategories(r)

	baseURL := helpers.GetFullRouteUrl(r)

	paginatedResponse := pagination.NewPaginatedResponse(categories.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, paginatedResponse)
	}
}

func (ch *CategoryHandlers) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	category, err := ch.Service.FindCategoryById(id)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, category.ToCategoryDTO())
	}
}

func (ch *CategoryHandlers) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryRequest dto.UpdateCategoryRequest
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Handle the error - the ID is not a valid integer
		helpers.WriteResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&categoryRequest)
	if err != nil {
		helpers.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateCategory(&categoryRequest); err != nil {
		helpers.WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	category, errCat := ch.Service.UpdateCategory(id, categoryRequest)
	if errCat != nil {
		helpers.WriteResponse(w, errCat.Code, errCat)
	} else {
		helpers.WriteResponse(w, http.StatusOK, category.ToCategoryDTO())
	}
}

func NewCategoryHandlers(service services.CategoryService) *CategoryHandlers {
	return &CategoryHandlers{
		Service: service,
	}
}
