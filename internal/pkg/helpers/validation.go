package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string
	Message string
}

func GetValidationErrors(err error) []ValidationError {
	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		var element ValidationError
		element.Field = err.Field()

		// Customize error messages based on the validation tag
		switch err.Tag() {
		case "required":
			element.Message = fmt.Sprintf("%s field is required", err.Field())
		case "min":
			element.Message = fmt.Sprintf("%s field must be at least %s characters long", err.Field(), err.Param())
		case "max":
			element.Message = fmt.Sprintf("%s field must not exceed %s characters", err.Field(), err.Param())
		case "lowercase":
			element.Message = fmt.Sprintf("%s field must be lowercase", err.Field())
		default:
			element.Message = fmt.Sprintf("%s field validation failed on %s", err.Field(), err.Tag())
		}

		errors = append(errors, element)
	}

	return errors
}
