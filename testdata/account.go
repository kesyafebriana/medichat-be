package testdata

import "medichat-be/domain"

var (
	AdminAccountNotVerified = domain.Account{
		ID:            1,
		Email:         "admin@example.com",
		EmailVerified: false,
		Role:          domain.AccountRoleAdmin,
		AccountType:   domain.AccountTypeRegular,
	}
	AdminAccount = domain.Account{
		ID:            1,
		Email:         "admin@example.com",
		EmailVerified: true,
		Role:          domain.AccountRoleAdmin,
		AccountType:   domain.AccountTypeRegular,
	}
	AdminPassword           = "123pass"
	AdminNewPassword        = "194nK#01204[Ma1!,"
	AdminHashedPassword     = "1088h8209jf09u1ih03h02nd"
	AdminNewHashedPassword  = "19471WRJ20g##G238g92F33$"
	AdminAccessToken        = "1093970735818098481"
	AdminRefreshToken       = "1093939529826982690"
	AdminResetPasswordToken = "10952789yhf09hn#JD38c#H3b8h2j987g2lkdj984t"
	AdminTokens             = domain.AuthTokens{
		AccessToken:  AdminAccessToken,
		RefreshToken: AdminRefreshToken,
	}

	AliceAccountNotVerified = domain.Account{
		ID:            1,
		Email:         "alice@example.com",
		EmailVerified: false,
		Role:          domain.AccountRoleUser,
		AccountType:   domain.AccountTypeRegular,
	}
	AliceAccount = domain.Account{
		ID:            1,
		Email:         "alice@example.com",
		EmailVerified: true,
		Role:          domain.AccountRoleUser,
		AccountType:   domain.AccountTypeRegular,
	}
	AlicePassword           = "123pass"
	AliceNewPassword        = "194nK#01204[Ma1!,"
	AliceHashedPassword     = "1088h8209jf09u1ih03h02nd"
	AliceNewHashedPassword  = "19471WRJ20g##G238g92F33$"
	AliceAccessToken        = "1093970735818098481"
	AliceRefreshToken       = "1093939529826982690"
	AliceResetPasswordToken = "10952789yhf09hn#JD38c#H3b8h2j987g2lkdj984t"
	AliceTokens             = domain.AuthTokens{
		AccessToken:  AliceAccessToken,
		RefreshToken: AliceRefreshToken,
	}

	DrBobAccountNotVerified = domain.Account{
		ID:            1,
		Email:         "dr.bob@example.com",
		EmailVerified: false,
		Role:          domain.AccountRoleDoctor,
		AccountType:   domain.AccountTypeRegular,
	}
	DrBobAccount = domain.Account{
		ID:            1,
		Email:         "dr.bob@example.com",
		EmailVerified: true,
		Role:          domain.AccountRoleDoctor,
		AccountType:   domain.AccountTypeRegular,
	}
	DrBobPassword           = "123pass"
	DrBobNewPassword        = "194nK#01204[Ma1!,"
	DrBobHashedPassword     = "1088h8209jf09u1ih03h02nd"
	DrBobNewHashedPassword  = "19471WRJ20g##G238g92F33$"
	DrBobAccessToken        = "1093970735818098481"
	DrBobRefreshToken       = "1093939529826982690"
	DrBobResetPasswordToken = "10952789yhf09hn#JD38c#H3b8h2j987g2lkdj984t"
	DrBobTokens             = domain.AuthTokens{
		AccessToken:  DrBobAccessToken,
		RefreshToken: DrBobRefreshToken,
	}

	PhBillAccountNotVerified = domain.Account{
		ID:            1,
		Email:         "ph.bill@example.com",
		EmailVerified: false,
		Role:          domain.AccountRolePharmacyManager,
		AccountType:   domain.AccountTypeRegular,
	}
	PhBillAccount = domain.Account{
		ID:            1,
		Email:         "ph.bill@example.com",
		EmailVerified: true,
		Role:          domain.AccountRolePharmacyManager,
		AccountType:   domain.AccountTypeRegular,
	}
	PhBillPassword           = "123pass"
	PhBillNewPassword        = "194nK#01204[Ma1!,"
	PhBillHashedPassword     = "1088h8209jf09u1ih03h02nd"
	PhBillNewHashedPassword  = "19471WRJ20g##G238g92F33$"
	PhBillAccessToken        = "1093970735818098481"
	PhBillRefreshToken       = "1093939529826982690"
	PhBillResetPasswordToken = "10952789yhf09hn#JD38c#H3b8h2j987g2lkdj984t"
	PhBillTokens             = domain.AuthTokens{
		AccessToken:  PhBillAccessToken,
		RefreshToken: PhBillRefreshToken,
	}
)
