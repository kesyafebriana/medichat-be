package dto

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MultipartForm[F any, D any] struct {
	Form F
	Data D
}

func ShouldBindMultipart[F any, D any](
	ctx *gin.Context,
	obj *MultipartForm[F, D],
) error {
	err := ctx.ShouldBind(&obj.Form)
	if err != nil {
		return err
	}

	data, ok := ctx.Request.MultipartForm.Value["data"]
	if !ok {
		return errors.New("data is required")
	}
	if len(data) == 0 {
		return errors.New("data is required")
	}

	return binding.JSON.BindBody([]byte(data[0]), &obj.Data)
}
