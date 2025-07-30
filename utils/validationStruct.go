package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	// Initialize the validator
	validate = validator.New()
}

func ValidateStruct(s any) error {
	// Validate the user struct
	if err := validate.Struct(s); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field()+" failed on the '"+err.Tag()+"' tag")
		}
		return errors.New("validation failed: " + strings.Join(validationErrors, ", "))
	}
	return nil
}
