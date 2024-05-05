package service

import (
	"errors"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/dto"
	"medichat-be/util"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type ChatService interface {
	PostMessage(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error)
	PostFile(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error)
	CreateRoom(req *dto.ChatRoom,ctx *gin.Context) (error)
	CloseRoom(roomId string,ctx *gin.Context) (error)
}

type chatService struct {
	dataRepository domain.DataRepository
	client  *firestore.Client
	cloud util.CloudinaryProvider
}

type ChatServiceOpts struct {
	DataRepository domain.DataRepository
	Client  *firestore.Client
	Cloud util.CloudinaryProvider
}
func NewChatService(opts ChatServiceOpts) *chatService {
	return &chatService{
		dataRepository: opts.DataRepository,
		client: opts.Client,
		cloud: opts.Cloud,
	}
}

func (u *chatService) CreateRoom(req *dto.ChatRoom,ctx *gin.Context) (error) {

	chatRepository := u.dataRepository.ChatRepository()


		room, err := chatRepository.AddRoom(ctx,
			req.UserId,
			req.DoctorId,
			req.End,
		)
		if err != nil{
			return err
		}

	colRef := u.client.Collection("rooms");

	roomId := strconv.Itoa(int(room.ID))
	_,err = colRef.Doc(roomId).Set(ctx,req)

	if err!= nil {
        return err
    }

	return nil

}

func (u *chatService) CloseRoom(roomId string,ctx *gin.Context) (error) {

	chatRepository := u.dataRepository.ChatRepository()

	colRef := u.client.Collection("rooms");

	_,err := colRef.Doc(roomId).Update(ctx,[]firestore.Update{
		{Path: "endAt", Value: time.Now()},
	})
	if err!= nil {
        return err
    }

	ss,err := colRef.Doc(roomId).Get(ctx)
	data := ss.Data()
	for i := 0; i < len(data); i++ {
		chat:= domain.Chat{
			ID: data["id"].(int64),
			RoomId: data["roomId"].(int64),
            Message: data["message"].(string),
            File: data["file"].(string),
            Type: data["type"].(string),
            UserId: data["userId"].(int),
			UserName: data["userName"].(string),
		}
		chatRepository.AddChat(ctx, chat)
	}
	if err!= nil {
        return err
    }

	return nil

}



func (u *chatService) PostMessage(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error) {

	colRef := u.client.Collection("rooms");
	content := map[string]interface{}{
        "userId": req.UserId,
        "userName": req.UserName,
        "message": req.Message,
		"url" : "",
        "createdAt": req.CreatedAt,
        "type": "message/text",
	}
	_,_,err := colRef.Doc(roomId).Collection("chats").Add(ctx, content)
	if err!= nil {
        return err
    }
	return nil

}

func (u *chatService) PostFile(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error) {

	fileType := req.File.Header.Get("Content-Type")

	if(req.File.Size >= constants.MaxFileSize){
		return errors.New("file size exceeded 5 Mb")
	}

	file,err := req.File.Open()
	if err!= nil {
        return err
    }

	fileName := req.File.Filename

	stringType := ""

	if (fileType == "application/pdf"){
		stringType = "message/pdf"
	}else if (fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/webp"){
		stringType = "message/image"
	}

	if(fileType == "application/pdf" || fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/webp"){
		opts := util.SendFileOpts{
			Context: ctx,
			Filename: fileName,
			Roomid: roomId,
			File: file,
		}
		response, err := u.cloud.SendFile(opts)

		if err != nil {
			return err
		}

		colRef := u.client.Collection("rooms");

		content := map[string]interface{}{
			"userId": req.UserId,
			"userName": req.UserName,
			"message": fileName,
			"url" : response.SecureURL,
			"createdAt": req.CreatedAt,
			"type": stringType,
		}
		_,_,err = colRef.Doc(roomId).Collection("chats").Add(ctx, content)
		if err!= nil {
			return err
		}
	} else{
		return errors.New("invalid file type")
	}

	return nil

}
