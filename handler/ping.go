package handler

import (
	"medichat-be/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Ping(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		dto.ResponseOk("pong"),
	)
}
