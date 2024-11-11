package dto

import (
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type LoginValidator interface {
	Validate() *helpers.ValidationResponse
}

type NewLoginRequest struct {
	Email    string `json:"email" validate:"required,email,min=3,max=250"`
	Password string `json:"password" validate:"required,min=3,max=250"`
}

func (req *NewLoginRequest) Validate() *helpers.ValidationResponse {
	return helpers.ValidateRequests(req)
}

// ValidateCategory is now a generic function that can handle any LoginValidator
func ValidateLogin(login LoginValidator) *helpers.ValidationResponse {
	return login.Validate()
}
