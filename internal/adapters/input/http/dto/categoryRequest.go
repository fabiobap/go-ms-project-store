package dto

import (
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type CategoryValidator interface {
	Validate() *helpers.ValidationResponse
}

type NewCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=250"`
	Slug string `json:"slug" validate:"omitempty,min=3,max=250"`
}

type CategoryUpdatedRequest struct {
	Name string `json:"name" validate:"required,min=3,max=250"`
	Slug string `json:"slug" validate:"omitempty,min=3,max=250"`
}

// Implement the Validate method for NewCategoryRequest
func (ncr *NewCategoryRequest) Validate() *helpers.ValidationResponse {
	return helpers.ValidateRequests(ncr)
}

// Implement the Validate method for CategoryUpdatedRequest
func (cur *CategoryUpdatedRequest) Validate() *helpers.ValidationResponse {
	return helpers.ValidateRequests(cur)
}

// ValidateCategory is now a generic function that can handle any CategoryValidator
func ValidateCategory(category CategoryValidator) *helpers.ValidationResponse {
	return category.Validate()
}
