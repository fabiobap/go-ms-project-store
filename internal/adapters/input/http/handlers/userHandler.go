package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-ms-project-store/internal/core/services"
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/go-ms-project-store/internal/pkg/pagination"
)

type UserHandlers struct {
	Service services.UserService
}

func (ch *UserHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := ch.Service.DeleteUser(id)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusNoContent, "")
	}
}

func (ch *UserHandlers) GetAllUserAdmins(w http.ResponseWriter, r *http.Request) {
	users, totalRows, filter, err := ch.Service.GetAllUserAdmins(r)

	baseURL := helpers.GetFullRouteUrl(r)

	paginatedResponse := pagination.NewPaginatedResponse(users.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, paginatedResponse)
	}
}

func (ch *UserHandlers) GetAllUserCustomers(w http.ResponseWriter, r *http.Request) {
	users, totalRows, filter, err := ch.Service.GetAllUserCustomers(r)

	baseURL := helpers.GetFullRouteUrl(r)

	paginatedResponse := pagination.NewPaginatedResponse(users.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, paginatedResponse)
	}
}

func (ch *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := ch.Service.FindUserById(id)
	if err != nil {
		helpers.WriteResponse(w, err.Code, err.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, user.ToUserDTO())
	}
}

func NewUserHandlers(service services.UserService) *UserHandlers {
	return &UserHandlers{
		Service: service,
	}
}
