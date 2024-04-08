package middleware

import (
	"medichat-be/apperror"
	"medichat-be/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors[0].Err

		apperr, ok := err.(*apperror.AppError)
		if ok {
			ctx.AbortWithStatusJSON(
				httpStatusFromAppError(apperr),
				dto.ResponseError(apperr),
			)
			return
		}

		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			dto.ResponseError(apperror.Wrap(err).(*apperror.AppError)),
		)
	}
}

func httpStatusFromAppError(err *apperror.AppError) int {
	switch err.Code {
	case apperror.CodeInternal:
		return http.StatusInternalServerError
	case apperror.CodeBadRequest:
		return http.StatusBadRequest
	case apperror.CodeValidationFailed:
		return http.StatusBadRequest
	case apperror.CodeConstraintViolation:
		return http.StatusBadRequest
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case apperror.CodeAlreadyExists:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperror.CodeInvalidToken:
		return http.StatusUnauthorized
	case apperror.CodeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
