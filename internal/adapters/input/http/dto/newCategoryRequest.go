package dto

import (
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/go-playground/validator/v10"
)

type NewCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=250"`
	Slug string `json:"slug" validate:"omitempty,min=3,max=250"`
}

func ValidateCategory(category *NewCategoryRequest) *helpers.ValidationResponse {
	validate := validator.New()
	err := validate.Struct(category)

	if err != nil {
		return helpers.GetValidationErrors(err)
	}

	return nil
}
