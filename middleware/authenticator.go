package middleware

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/cryptoutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticator(jwtProvider cryptoutil.JWTProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authVal := ctx.GetHeader("Authorization")

		if authVal == "" {
			ctx.Error(apperror.NewUnauthorized(nil))
			ctx.Abort()
			return
		}

		tokens := strings.Fields(authVal)
		if len(tokens) != 2 {
			ctx.Error(apperror.NewInvalidToken(nil))
			ctx.Abort()
			return
		}
		if tokens[0] != "Bearer" {
			ctx.Error(apperror.NewInvalidToken(nil))
			ctx.Abort()
			return
		}

		token := tokens[1]

		claims, err := jwtProvider.VerifyToken(token)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(constants.ContextAccountID, claims.UserID)

		ctx.Next()
	}
}
