package app

import (
	"net/http"

	"github.com/go-ms-project-store/app/handlers"
	"github.com/go-ms-project-store/service"
)

type CategoryHandlers struct {
	service service.CategoryService
}

func (ch *CategoryHandlers) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := ch.service.GetAllCategories()
	if err != nil {
		handlers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		handlers.WriteResponse(w, http.StatusOK, categories)
	}
}
