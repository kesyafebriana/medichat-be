package dto

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MultipartForm[T any] struct {
	Files []*multipart.FileHeader `form:"file"`
	Data  T                       `form:"data" binding:"required"`
}

type multipartFormInternal struct {
	Files []*multipart.FileHeader `form:"file"`
	Data  string                  `form:"data" binding:"required"`
}

func ShouldBindMultipart[T any](
	ctx *gin.Context,
	form *MultipartForm[T],
) error {
	var tmp multipartFormInternal

	err := ctx.ShouldBind(&tmp)
	if err != nil {
		return err
	}

	form.Files = tmp.Files

	return binding.JSON.BindBody([]byte(tmp.Data), &form.Data)
}
