package service_test

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/mocks/cryptomocks"
	"medichat-be/mocks/domainmocks"
	"medichat-be/service"
	"medichat-be/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			ctx := context.Background()

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

			refreshProv.On(
				"CreateToken",
				tt.getAccountWithCreds.Val.Account.ID,
			).Return(
				tt.createRefreshToken.Val,
				tt.createRefreshToken.Err,
			)

			opts := service.AccountServiceOpts{
				DataRepository: dataRepo,
				PasswordHasher: pwdHasher,
			}

			switch tt.getAccountWithCreds.Val.Account.Role {
			case domain.AccountRoleAdmin:
				opts.AdminAccessProvider = accessProv
				opts.AdminRefreshProvider = refreshProv
			case domain.AccountRoleUser:
				opts.UserAccessProvider = accessProv
				opts.UserRefreshProvider = refreshProv
			case domain.AccountRoleDoctor:
				opts.DoctorAccessProvider = accessProv
				opts.DoctorRefreshProvider = refreshProv
			case domain.AccountRolePharmacyManager:
				opts.PharmacyManagerAccessProvider = accessProv
				opts.PharmacyManagerRefreshProvider = refreshProv
			}

			s := service.NewAccountService(opts)

			// when
			got, err := s.Login(ctx, tt.creds)

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
