package constants

import "medichat-be/domain"

var (
	AvailableAccountRoles = map[string]bool{
		domain.AccountRoleAdmin:           true,
		domain.AccountRoleUser:            true,
		domain.AccountRoleDoctor:          true,
		domain.AccountRolePharmacyManager: true,
	}
)

const (
	PasswordMinLength         = 8
	PasswordMaxLength         = 24
	PasswordSpecialCharacters = "!@#$%^&*()\\-_=+{};:,<.>]"
)
