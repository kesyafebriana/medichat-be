package middleware

import (
	"errors"
	"fmt"
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(logger logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		method := ctx.Request.Method

		ctx.Next()

		if query != "" {
			path = path + "?" + query
		}

		statusCode := ctx.Writer.Status()

		requestID, exists := ctx.Get(constants.ContextRequestID)
		if !exists {
			requestID = ""
		}

		fields := map[string]interface{}{
			"request_id":  requestID,
			"path":        path,
			"latency":     time.Since(start),
			"method":      method,
			"status_code": statusCode,
		}

		if statusCode >= 500 && statusCode < 600 {
			var appErr *apperror.AppError
			for _, err := range ctx.Errors {
				fields["error"] = err
				fields["error_type"] = fmt.Sprintf("%T", err)
				logger.ErrorFields(fields, "internal error")
				if errors.As(err, &appErr) {
					logger.Errorf("Stack trace:\n %s", string(appErr.GetStackTrace()))
				}
			}
			return
		}

		if statusCode >= 400 && statusCode < 500 {
			for _, err := range ctx.Errors {
				fields["error"] = err
				fields["error_type"] = fmt.Sprintf("%T", err)
				logger.InfoFields(fields, "request processed with user error")
			}
			return
		}
		logger.InfoFields(fields, "request processed")
	}
}
