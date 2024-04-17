package service_test

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/cryptoutil"
	"medichat-be/domain"
	"medichat-be/mocks/cryptomocks"
	"medichat-be/mocks/domainmocks"
	"medichat-be/service"
	"medichat-be/testdata"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func Test_accountService_Register(t *testing.T) {
	tests := []struct {
		name string

		isAccountExist testdata.Result[bool]
		addAccount     testdata.Result[domain.Account]

		ctx   context.Context
		creds domain.AccountRegisterCredentials

		want testdata.WantValue[domain.Account]
	}{
		{
			name: "should successfully register alice",

			isAccountExist: testdata.Result[bool]{
				Val: false,
			},
			addAccount: testdata.Result[domain.Account]{
				Val: testdata.AliceAccount,
			},

			ctx: context.Background(),
			creds: domain.AccountRegisterCredentials{
				Account:  testdata.AliceAccount,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.Account]{
				Val: testdata.AliceAccount,
			},
		},
		{
			name: "should return already exists when email already used",

			isAccountExist: testdata.Result[bool]{
				Val: true,
			},

			ctx: context.Background(),
			creds: domain.AccountRegisterCredentials{
				Account:  testdata.AliceAccount,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.Account]{
				Err: apperror.CodeAlreadyExists,
			},
		},
		{
			name: "should return internal error when check email exists",

			isAccountExist: testdata.Result[bool]{
				Err: apperror.NewInternal(nil),
			},

			ctx: context.Background(),
			creds: domain.AccountRegisterCredentials{
				Account:  testdata.AliceAccount,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.Account]{
				Err: apperror.CodeInternal,
			},
		},
		{
			name: "should return internal error when add account",

			isAccountExist: testdata.Result[bool]{
				Val: false,
			},
			addAccount: testdata.Result[domain.Account]{
				Err: apperror.NewInternal(nil),
			},

			ctx: context.Background(),
			creds: domain.AccountRegisterCredentials{
				Account:  testdata.AliceAccount,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.Account]{
				Err: apperror.CodeInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			accountRepo := new(domainmocks.AccountRepository)
			dataRepo := testdata.NewDataRepositoryMock(testdata.DataRepositoryMockOpts{
				AccountRepository: accountRepo,
			})
			pwdHasher := new(cryptomocks.PasswordHasher)

			accountRepo.On(
				"IsExistByEmail",
				tt.ctx,
				tt.creds.Account.Email,
			).Return(
				tt.isAccountExist.Val,
				tt.isAccountExist.Err,
			)

			accountRepo.On(
				"Add",
				tt.ctx,
				domain.AccountWithCredentials{
					Account: tt.creds.Account,
				},
			).Return(
				tt.addAccount.Val,
				tt.addAccount.Err,
			)

			opts := service.AccountServiceOpts{
				DataRepository: dataRepo,
				PasswordHasher: pwdHasher,
			}

			s := service.NewAccountService(opts)

			testdata.OnDataRepositoryAtomic(
				dataRepo,
				tt.ctx,
				s.RegisterClosure(tt.ctx, tt.creds),
			)

			// when
			got, err := s.Register(tt.ctx, tt.creds)

			// then
			assert.Equal(t, tt.want.Val, got)
			if tt.want.Err != 0 {
				apperror.AssertErrorIsCode(t, err, tt.want.Err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_accountService_Login(t *testing.T) {
	tests := []struct {
		name string

		getAccountWithCreds testdata.Result[domain.AccountWithCredentials]
		checkPwdErr         error
		createAccessToken   testdata.Result[string]
		createRefreshToken  testdata.Result[string]

		ctx   context.Context
		creds domain.AccountLoginCredentials

		want testdata.WantValue[domain.AuthTokens]
	}{
		{
			name: "should successfully login as Alice",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Val: domain.AccountWithCredentials{
					Account:        testdata.AliceAccount,
					HashedPassword: testdata.AliceHashedPassword,
				},
			},
			checkPwdErr: nil,
			createAccessToken: testdata.Result[string]{
				Val: testdata.AliceAccessToken,
			},
			createRefreshToken: testdata.Result[string]{
				Val: testdata.AliceRefreshToken,
			},

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.AliceAccount.Email,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Val: testdata.AliceTokens,
			},
		},
		{
			name: "should successfully login as Admin",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Val: domain.AccountWithCredentials{
					Account:        testdata.AdminAccount,
					HashedPassword: testdata.AdminHashedPassword,
				},
			},
			checkPwdErr: nil,
			createAccessToken: testdata.Result[string]{
				Val: testdata.AdminAccessToken,
			},
			createRefreshToken: testdata.Result[string]{
				Val: testdata.AdminRefreshToken,
			},

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.AdminAccount.Email,
				Password: testdata.AdminPassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Val: testdata.AdminTokens,
			},
		},
		{
			name: "should successfully login as doctor Bob",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Val: domain.AccountWithCredentials{
					Account:        testdata.DrBobAccount,
					HashedPassword: testdata.DrBobHashedPassword,
				},
			},
			checkPwdErr: nil,
			createAccessToken: testdata.Result[string]{
				Val: testdata.DrBobAccessToken,
			},
			createRefreshToken: testdata.Result[string]{
				Val: testdata.DrBobRefreshToken,
			},

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.DrBobAccount.Email,
				Password: testdata.DrBobPassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Val: testdata.DrBobTokens,
			},
		},
		{
			name: "should successfully login as pharmacist Bill",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Val: domain.AccountWithCredentials{
					Account:        testdata.PhBillAccount,
					HashedPassword: testdata.PhBillHashedPassword,
				},
			},
			checkPwdErr: nil,
			createAccessToken: testdata.Result[string]{
				Val: testdata.PhBillAccessToken,
			},
			createRefreshToken: testdata.Result[string]{
				Val: testdata.PhBillRefreshToken,
			},

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.PhBillAccount.Email,
				Password: testdata.PhBillPassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Val: testdata.PhBillTokens,
			},
		},
		{
			name: "should return user not found",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Err: apperror.NewNotFound(),
			},

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.AliceAccount.Email,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Err: apperror.CodeNotFound,
			},
		},
		{
			name: "should return unauthoried when wrong password",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Val: domain.AccountWithCredentials{
					Account:        testdata.AliceAccount,
					HashedPassword: testdata.AliceHashedPassword,
				},
			},
			checkPwdErr: apperror.NewWrongPassword(nil),

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.AliceAccount.Email,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Err: apperror.CodeUnauthorized,
			},
		},
		{
			name: "should return internal error when check password",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Val: domain.AccountWithCredentials{
					Account:        testdata.AliceAccount,
					HashedPassword: testdata.AliceHashedPassword,
				},
			},
			checkPwdErr: apperror.NewInternal(nil),

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.AliceAccount.Email,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Err: apperror.CodeInternal,
			},
		},
		{
			name: "should return internal error when generating access token",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Val: domain.AccountWithCredentials{
					Account:        testdata.AliceAccount,
					HashedPassword: testdata.AliceHashedPassword,
				},
			},
			checkPwdErr: nil,
			createAccessToken: testdata.Result[string]{
				Err: apperror.NewInternal(nil),
			},

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.AliceAccount.Email,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Err: apperror.CodeInternal,
			},
		},
		{
			name: "should return internal error when generating refresh token",

			getAccountWithCreds: testdata.Result[domain.AccountWithCredentials]{
				Val: domain.AccountWithCredentials{
					Account:        testdata.AliceAccount,
					HashedPassword: testdata.AliceHashedPassword,
				},
			},
			checkPwdErr: nil,
			createAccessToken: testdata.Result[string]{
				Val: testdata.AliceAccessToken,
			},
			createRefreshToken: testdata.Result[string]{
				Err: apperror.NewInternal(nil),
			},

			ctx: context.Background(),
			creds: domain.AccountLoginCredentials{
				Email:    testdata.AliceAccount.Email,
				Password: testdata.AlicePassword,
			},

			want: testdata.WantValue[domain.AuthTokens]{
				Err: apperror.CodeInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			accountRepo := new(domainmocks.AccountRepository)
			dataRepo := testdata.NewDataRepositoryMock(testdata.DataRepositoryMockOpts{
				AccountRepository: accountRepo,
			})
			pwdHasher := new(cryptomocks.PasswordHasher)
			accessProv := new(cryptomocks.JWTProvider)
			refreshProv := new(cryptomocks.JWTProvider)

			accountRepo.On(
				"GetWithCredentialsByEmail",
				tt.ctx,
				tt.creds.Email,
			).Return(
				tt.getAccountWithCreds.Val,
				tt.getAccountWithCreds.Err,
			)

			pwdHasher.On(
				"CheckPassword",
				tt.getAccountWithCreds.Val.HashedPassword,
				tt.creds.Password,
			).Return(
				tt.checkPwdErr,
			)

			accessProv.On(
				"CreateToken",
				tt.getAccountWithCreds.Val.Account.ID,
			).Return(
				tt.createAccessToken.Val,
				tt.createAccessToken.Err,
			)

			accessProv.On(
				"VerifyToken",
				tt.createAccessToken.Val,
			).Return(
				cryptoutil.JWTClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Time{}),
					},
				},
				nil,
			)

			refreshProv.On(
				"CreateToken",
				tt.getAccountWithCreds.Val.Account.ID,
			).Return(
				tt.createRefreshToken.Val,
				tt.createRefreshToken.Err,
			)

			refreshProv.On(
				"VerifyToken",
				tt.createRefreshToken.Val,
			).Return(
				cryptoutil.JWTClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Time{}),
					},
				},
				nil,
			)

			opts := service.AccountServiceOpts{
				DataRepository:  dataRepo,
				PasswordHasher:  pwdHasher,
				RefreshProvider: refreshProv,
			}

			switch tt.getAccountWithCreds.Val.Account.Role {
			case domain.AccountRoleAdmin:
				opts.AdminAccessProvider = accessProv
			case domain.AccountRoleUser:
				opts.UserAccessProvider = accessProv
			case domain.AccountRoleDoctor:
				opts.DoctorAccessProvider = accessProv
			case domain.AccountRolePharmacyManager:
				opts.PharmacyManagerAccessProvider = accessProv
			}

			s := service.NewAccountService(opts)

			// when
			got, err := s.Login(tt.ctx, tt.creds)

			// then
			assert.Equal(t, tt.want.Val, got)
			if tt.want.Err != 0 {
				apperror.AssertErrorIsCode(t, err, tt.want.Err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
