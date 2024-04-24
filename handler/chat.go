package handler

import (
	"medichat-be/apperror"
	"medichat-be/constants"
	"medichat-be/dto"
	"medichat-be/service"
	"net/http"
	"strconv"
	"time"

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
	var req dto.ChatMessage
	req.Type = ctx.PostForm("type")
	req.UserName = ctx.PostForm("userName")

	timeChat,err := time.Parse("2006-01-02T15:04:05Z07:00",ctx.PostForm("createdAt"))
	if err != nil {
		ctx.Error(err)
        ctx.Abort()
        return
	}
	req.CreatedAt = timeChat
	userId,err := strconv.Atoi(ctx.PostForm("userId"))
	if err!= nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}
	req.UserId = userId
	if roomId == "" {
		ctx.Error(apperror.NewBadRequest(nil))
		ctx.Abort()
		return
	}

	// Implement Static Type !!!
	if (req.Type == "text"){
		req.Message = ctx.PostForm("message")
		if req.Message == "" {
			ctx.Error(apperror.NewBadRequest(nil))
			ctx.Abort()
			return
		}
		err = h.chatService.PostMessage(&req,roomId,ctx)
	} else if(req.Type == "files"){
		fileHeader,err := ctx.FormFile("file")

		
		if err!= nil {
			ctx.Error(apperror.NewBadRequest(err))
            ctx.Abort()
            return
        }
		
		
		req.File= fileHeader

		err = h.chatService.PostFile(&req,roomId,ctx)
		if err!= nil {
            ctx.Error(apperror.NewBadRequest(err))
            ctx.Abort()
            return
        }
	}


	if err!= nil {
        ctx.Error(apperror.NewInternal(err))
        ctx.Abort()
        return
    }
	ctx.JSON(http.StatusOK, gin.H{"message": "message sent"})

}

func (h *ChatHandler) CreateRoom(ctx *gin.Context) {
	var req dto.ChatRoom

	req.UserName = ctx.PostForm("userName")
	userId,err := strconv.Atoi(ctx.PostForm("userId"))
	if err != nil {
		ctx.Error(apperror.NewBadRequest(err))
        ctx.Abort()
        return
	}
	req.UserId = userId

	req.DoctorName = ctx.PostForm("doctorName")
	doctorId,err := strconv.Atoi(ctx.PostForm("doctorId"))
	if err != nil{
		ctx.Error(apperror.NewBadRequest(err))
        ctx.Abort()
        return
	}
	req.DoctorId = doctorId

	date,err := time.Parse("2006-01-02T15:04:05Z07:00",ctx.PostForm("date"))
	if err != nil {
		ctx.Error(err)
        ctx.Abort()
        return
	}
	req.Start = date

	extra , _ := time.ParseDuration(constants.ChatDuration)

	req.End = date.Add(extra)

	err = h.chatService.CreateRoom(&req,ctx)
	if err!= nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "message sent"})

}

func (h *ChatHandler) CloseRoom(ctx *gin.Context) {

	roomId := ctx.Query("roomId")

	err := h.chatService.CloseRoom(roomId,ctx)
	if err!= nil {
		ctx.Error(apperror.NewBadRequest(err))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "message sent"})

}
