package handler

import (
	"encoding/json"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) Chat(ctx *gin.Context) {

	roomId := ctx.Query("roomId")
	var req domain.ChatMessage

	err := json.NewDecoder(ctx.Request.Body).Decode(&req)
	

	if err!= nil {
		ctx.Error(apperror.NewBadRequest(err))
        ctx.Abort()
        return
	}

	if roomId == "" {
        ctx.Error(apperror.NewBadRequest(nil))
        ctx.Abort()
        return
    }

	if req.Message == "" {
        ctx.Error(apperror.NewBadRequest(nil))
        ctx.Abort()
        return
    }

	err = h.chatService.Post(&req,roomId,ctx)
	if err!= nil {
        ctx.Error(apperror.NewInternal(err))
        ctx.Abort()
        return
    }
	ctx.JSON(http.StatusOK, gin.H{"message": "message sent"})

}
