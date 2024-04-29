package util

import (
	"fmt"
	"log"
	"medichat-be/constants"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LimitContentLength(ctx *gin.Context, max_length int) error {
	lengthStr := ctx.Request.Header.Get(constants.ContentLength)

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return err
	}

	log.Println(length)

	if length > max_length {
		return fmt.Errorf("content length exceeds %d bytes (got %d bytes)", max_length, length)
	}

	return nil
}
