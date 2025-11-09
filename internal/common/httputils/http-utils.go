package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func DecodeAndValidateRequest(r *http.Request, dest interface{}, v *validator.Validate) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return fmt.Errorf("request body cannot be empty")
	}

	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := v.Struct(dest); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("validation setup error: %w", err)
		}

		var errMsg string
		for _, e := range err.(validator.ValidationErrors) {
			errMsg += fmt.Sprintf("Field '%s' failed validation tag '%s'; ", e.Field(), e.Tag())
		}
		return fmt.Errorf("validation failed: %s", errMsg)
	}

	return nil
}

func ConvertValidationErrors(err error) []FieldError {
	var fieldErrors []FieldError

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   e.Field(),
				Message: "failed on the '" + e.Tag() + "' rule",
			})
		}
	}

	return fieldErrors
}
