package service

import (
	"errors"
	"fmt"
	"medichat-be/constants"
	"medichat-be/dto"
	"medichat-be/util"

	"cloud.google.com/go/firestore"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type ChatService interface {
	PostMessage(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error)
	PostFile(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error)
	CreateRoom(req *dto.ChatRoom,ctx *gin.Context) (error)
	CloseRoom(roomId string,ctx *gin.Context) (error)
}

type chatServiceImpl struct {
	client *firestore.Client
	cloud util.CloudinaryProvider

}

func NewChatServiceImpl(client  *firestore.Client, cloud util.CloudinaryProvider) *chatServiceImpl {
	return &chatServiceImpl{
		client: client,
		cloud: cloud,
	}
}

func (u *chatServiceImpl) CreateRoom(req *dto.ChatRoom,ctx *gin.Context) (error) {

	colRef := u.client.Collection("rooms");
	_, _,err := colRef.Add(ctx,req)

	if err!= nil {
        return err
    }

	return nil

}

func (u *chatServiceImpl) CloseRoom(roomId string,ctx *gin.Context) (error) {

	colRef := u.client.Collection("rooms");
	_,err := colRef.Doc(roomId).Update(ctx,[]firestore.Update{
		{Path: "open", Value: false},
	})

	if err!= nil {
        return err
    }

	return nil

}



func (u *chatServiceImpl) PostMessage(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error) {


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

func (u *chatServiceImpl) PostFile(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error) {

	fileType := req.File.Header.Get("Content-Type")

	if(req.File.Size >= constants.MaxFileSize){
		return errors.New("file size exceeded 5 Mb")
	}

	file,err := req.File.Open()
	if err!= nil {
        return err
    }

	fileName := req.File.Filename

	resp := make(chan *uploader.UploadResult)
	stringType := ""

	if (fileType == "application/pdf"){
		stringType = "message/pdf"
	}else if (fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/webp"){
		stringType = "message/image"
	}

	

	if(fileType == "application/pdf" || fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/webp"){
		go func() {
			opts := util.SendFileOpts{
				Context: ctx,
				Filename: fileName,
				Roomid: roomId,
				File: file,
			}
			res, err := u.cloud.SendFile(opts)
			if err != nil {
				fmt.Println(err)
				return
			}
			resp<- res
		}()

		response := <-resp
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
