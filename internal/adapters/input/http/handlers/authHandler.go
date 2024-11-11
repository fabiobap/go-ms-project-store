package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/core/services"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type AuthHandlers struct {
	Service services.AuthService
}

func (ah *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.NewLoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		helpers.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateLogin(&loginRequest); err != nil {
		helpers.WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenRes, errT := ah.Service.Login(loginRequest)
	if errT != nil {
		helpers.WriteResponse(w, errT.Code, errT.AsMessage())
	} else {
		helpers.WriteResponse(w, http.StatusOK, tokenRes)
	}
}

// func (ch *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
// 	// paginatedResponse := pagination.NewPaginatedResponse(users.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
// 	if err != nil {
// 		helpers.WriteResponse(w, err.Code, err.AsMessage())
// 	} else {
// 		helpers.WriteResponse(w, http.StatusOK, paginatedResponse)
// 	}
// }

// func (ch *AuthHandlers) Me(w http.ResponseWriter, r *http.Request) {
// 	// paginatedResponse := pagination.NewPaginatedResponse(users.ToDTO(), filter.Page, filter.PerPage, int(totalRows), baseURL)
// 	if err != nil {
// 		helpers.WriteResponse(w, err.Code, err.AsMessage())
// 	} else {
// 		helpers.WriteResponse(w, http.StatusOK, paginatedResponse)
// 	}
// }

// func (ch *AuthHandlers) Refresh(w http.ResponseWriter, r *http.Request) {
// 	// user, err := ch.Service.FindAuthById(id)
// 	if err != nil {
// 		helpers.WriteResponse(w, err.Code, err.AsMessage())
// 	} else {
// 		helpers.WriteResponse(w, http.StatusOK, user.ToAuthDTO())
// 	}
// }

// func (ch *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
// 	// user, err := ch.Service.FindAuthById(id)
// 	if err != nil {
// 		helpers.WriteResponse(w, err.Code, err.AsMessage())
// 	} else {
// 		helpers.WriteResponse(w, http.StatusOK, user.ToAuthDTO())
// 	}
// }

func NewAuthHandlers(service services.AuthService) *AuthHandlers {
	return &AuthHandlers{
		Service: service,
	}
}
