package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func ValidateRequests(s interface{}) *ValidationResponse {
	validate := validator.New()
	err := validate.Struct(s)

	if err != nil {
		return getValidationErrors(err)
	}

	return nil
}

func getValidationErrors(err error) *ValidationResponse {
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string][]string)
		var firstMessage string

		for _, err := range validationErrors {
			field := err.Field()
			field = makeFirstLetterLower(field) // convert Name to name
			var message string

			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("The %s field is required.", field)
			case "min":
				message = fmt.Sprintf("The %s must be at least %s characters.", field, err.Param())
			case "max":
				message = fmt.Sprintf("The %s must not exceed %s characters.", field, err.Param())
			case "lowercase":
				message = fmt.Sprintf("The %s must be lowercase.", field)
			case "email":
				message = fmt.Sprintf("The %s must be a valid email address.", field)
			case "len":
				message = fmt.Sprintf("The %s must be exactly %s characters.", field, err.Param())
			case "uuid4":
				message = fmt.Sprintf("The %s must be a valid UUID v4.", field)
			case "dive":
				message = fmt.Sprintf("The %s contains invalid items.", field)
			case "numeric":
				message = fmt.Sprintf("The %s must contain only numbers.", field)
			case "gt":
				message = fmt.Sprintf("The %s must be greater than %s.", field, err.Param())
			case "gte":
				message = fmt.Sprintf("The %s must be greater than or equal to %s.", field, err.Param())
			case "lt":
				message = fmt.Sprintf("The %s must be less than %s.", field, err.Param())
			case "lte":
				message = fmt.Sprintf("The %s must be less than or equal to %s.", field, err.Param())
			default:
				message = fmt.Sprintf("The %s field is invalid.", field)
			}

			if firstMessage == "" {
				firstMessage = message
			}

			if errors[field] == nil {
				errors[field] = []string{}
			}
			errors[field] = append(errors[field], message)
		}

		return &ValidationResponse{
			Message: firstMessage,
			Errors:  errors,
		}
	}

	return nil
}

// Helper function to convert first letter to lowercase
func makeFirstLetterLower(s string) string {
	if len(s) == 0 {
		return s
	}
	return fmt.Sprintf("%c%s", s[0]|32, s[1:]) // convert first letter to lowercase
}
