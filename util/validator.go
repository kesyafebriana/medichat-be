package util

import (
	"medichat-be/constants"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("account_role", AccountRoleValidator)
		v.RegisterValidation("password", PasswordValidator)
		v.RegisterValidation("no_leading_trailing_space", NoLeadingOrTrailingSpaceValidator)
	}
}

func AccountRoleValidator(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	_, ok := constants.AvailableAccountRoles[s]

	return ok
}

func PasswordValidator(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	if len(s) < constants.PasswordMinLength ||
		len(s) > constants.PasswordMaxLength {
		return false
	}

	var countLower, countUpper, countNumber, countSpecial int

	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			countLower++
		} else if r >= 'A' && r <= 'Z' {
			countUpper++
		} else if r >= '0' && r <= '9' {
			countNumber++
		} else if strings.ContainsRune(constants.PasswordSpecialCharacters, r) {
			countSpecial++
		} else {
			return false
		}
	}

	return countLower > 0 && countUpper > 0 && countNumber > 0 && countSpecial > 0
}

func NoLeadingOrTrailingSpaceValidator(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	return strings.TrimSpace(s) == s
}
