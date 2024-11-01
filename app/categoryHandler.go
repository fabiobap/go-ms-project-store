package app

import (
	"net/http"

	"github.com/go-ms-project-store/service"
)

type CategoryHandlers struct {
	service service.CategoryService
}

func (ch *CategoryHandlers) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := ch.service.GetAllCategories()
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, categories)
	}
}
