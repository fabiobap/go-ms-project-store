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

type ProductHandlers struct {
	Service services.ProductService
}

func (ch *ProductHandlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	_, err := ch.Service.DeleteProduct(id)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusNoContent, "")
	}
}

func (ch *ProductHandlers) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productRequest dto.NewProductRequest

	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		helpers.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateProduct(&productRequest); err != nil {
		helpers.WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	product, errCat := ch.Service.CreateProduct(productRequest)
	if errCat != nil {
		helpers.WriteResponse(w, errCat.Code, errCat)
	} else {
		helpers.WriteResponse(w, http.StatusCreated, product.ToProductDTO())
	}
}

func (ch *ProductHandlers) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, totalRows, filter, err := ch.Service.GetAllProducts(r)

	baseURL := helpers.GetFullRouteUrl(r)

	paginatedResponse := pagination.NewPaginatedResponse(products.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, paginatedResponse)
	}
}

func (ch *ProductHandlers) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	product, err := ch.Service.FindProductById(id)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, product.ToProductDTO())
	}
}

func (ch *ProductHandlers) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var productRequest dto.UpdateProductRequest
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Handle the error - the ID is not a valid integer
		helpers.WriteResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		helpers.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateProduct(&productRequest); err != nil {
		helpers.WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	product, errCat := ch.Service.UpdateProduct(id, productRequest)
	if errCat != nil {
		helpers.WriteResponse(w, errCat.Code, errCat)
	} else {
		helpers.WriteResponse(w, http.StatusOK, product.ToProductDTO())
	}
}

func NewProductHandlers(service services.ProductService) *ProductHandlers {
	return &ProductHandlers{
		Service: service,
	}
}
