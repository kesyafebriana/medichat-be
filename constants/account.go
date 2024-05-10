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
	DefaultPhotoURL           = "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_960_720.png"
)
