package dto

import (
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type UserRegisterValidator interface {
	Validate() *helpers.ValidationResponse
}

type NewUserRegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=250"`
	Email                string `json:"email" validate:"required,email,max=250"`
	Password             string `json:"password" validate:"required,min=8,max=250"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=250"`
}

func (ncr *NewUserRegisterRequest) Validate() *helpers.ValidationResponse {
	if err := helpers.ValidateRequests(ncr); err != nil {
		return err
	}

	// Add password confirmation validation
	if ncr.Password != ncr.PasswordConfirmation {
		msg := make(map[string][]string)
		msg["password_confirmation"] = append(msg["password_confirmation"], "Password confirmation must match password")
		return &helpers.ValidationResponse{
			Message: "Password confirmation does not match",
			Errors:  msg,
		}
	}

	return nil
}

// ValidateUserRegister is now a generic function that can handle any UserRegisterValidator
func ValidateUserRegister(userRegister UserRegisterValidator) *helpers.ValidationResponse {
	return userRegister.Validate()
}
