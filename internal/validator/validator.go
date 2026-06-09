package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator for reusable validation across handlers.
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator instance.
func New() *Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
	return &Validator{validate: v}
}

// ValidateStruct validates a struct and returns field-level error details.
func (v *Validator) ValidateStruct(s interface{}) map[string]string {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	details := make(map[string]string)
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		details["_error"] = err.Error()
		return details
	}

	for _, fe := range validationErrors {
		field := fe.Field()
		details[field] = formatValidationError(fe)
	}
	return details
}

func formatValidationError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	default:
		return fmt.Sprintf("%s failed validation on '%s'", fe.Field(), fe.Tag())
	}
}
