package validation

import (
	"user-management/internal/user"

	"github.com/go-playground/validator/v10"
)

func RegisterValidations(validate *validator.Validate) {
	validate.RegisterValidation("userStatus", validateUserStatus)
}

func validateUserStatus(fl validator.FieldLevel) bool {
	statusField := fl.Field()

	userStatus, ok := statusField.Interface().(user.UserStatus)
	if !ok {
		return false
	}

	statusStr := userStatus.String()
	_, err := user.ParseUserStatus(statusStr)
	return err == nil
}
