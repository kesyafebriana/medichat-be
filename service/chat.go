package service

import (
	"fmt"
	"medichat-be/domain"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type ChatService interface {
	PostMessage(req *domain.ChatMessage,roomId string,ctx *gin.Context) (error)
	PostFile(fileType string,req *domain.ChatMessage,roomId string,ctx *gin.Context) (error)
}

type chatServiceImpl struct {
	client *firestore.Client
	cloud *cloudinary.Cloudinary

}

func NewChatServiceImpl(client  *firestore.Client, cloud *cloudinary.Cloudinary) *chatServiceImpl {
	return &chatServiceImpl{
		client: client,
		cloud: cloud,
	}
}

func (u *chatServiceImpl) PostMessage(req *domain.ChatMessage,roomId string,ctx *gin.Context) (error) {


	colRef := u.client.Collection("rooms");
	content := map[string]interface{}{
        "userId": req.UserId,
        "userName": req.UserName,
        "message": req.Message,
        "createdAt": req.CreatedAt,
        "type": "message/text",
	}
	_,_,err := colRef.Doc(roomId).Collection("chats").Add(ctx, content)
	if err!= nil {
        return err
    }
	return nil

}

func (u *chatServiceImpl) PostFile(fileType string,req *domain.ChatMessage,roomId string,ctx *gin.Context) (error) {

	fmt.Println("Upload File")

	now := time.Now()
	resp, err := u.cloud.Upload.Upload(ctx,*req.File,uploader.UploadParams{
		DisplayName: roomId+now.Format("2006_01_02_T15:04:05"),
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(true),
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
		
	}
	stringType := ""
	if (fileType == "application/pdf"){
		stringType = "message/pdf"
	}else if (fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/webp"){
		stringType = "message/image"
	}
	colRef := u.client.Collection("rooms");
	content := map[string]interface{}{
        "userId": req.UserId,
        "userName": req.UserName,
        "message": resp.SecureURL,
        "createdAt": req.CreatedAt,
        "type": stringType,
	}
	_,_,err = colRef.Doc(roomId).Collection("chats").Add(ctx, content)
	if err!= nil {
        return err
    }
	return nil

}
