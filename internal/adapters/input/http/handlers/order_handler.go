package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/adapters/input/http/middlewares"
	"github.com/go-ms-project-store/internal/core/ports"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type OrderHandlers struct {
	Service ports.OrderService
}

func (oh *OrderHandlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	user_id, ok := middlewares.GetUserID(r.Context())
	if !ok {
		helpers.WriteResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var orderRequest dto.NewOrderRequest

	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		helpers.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateOrder(&orderRequest); err != nil {
		helpers.WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	order, errCat := oh.Service.CreateOrder(orderRequest, user_id)
	if errCat != nil {
		helpers.WriteResponse(w, errCat.Code, errCat)
	} else {
		helpers.WriteResponse(w, http.StatusCreated, order.ToOrderDTO())
	}
}

func NewOrderHandlers(service ports.OrderService) *OrderHandlers {
	return &OrderHandlers{
		Service: service,
	}
}
