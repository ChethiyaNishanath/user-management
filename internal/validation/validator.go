package validation

import (
	"user-management/internal/user"

	"github.com/go-playground/validator/v10"
)

func RegisterValidations(validate *validator.Validate) {
	validate.RegisterValidation("userStatus", validateUserStatus)
}

func validateUserStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	_, err := user.ParseUserStatus(status)
	return err == nil
}
