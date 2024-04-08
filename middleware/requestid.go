package middleware

import (
	"medichat-be/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.NewRandom()
		if err == nil {
			idstr := id.String()
			ctx.Set(constants.ContextRequestID, idstr)
		}

		ctx.Next()
	}
}
