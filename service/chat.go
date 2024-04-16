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
	PostFile(req *domain.ChatMessage,roomId string,ctx *gin.Context) (error)
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

func (u *chatServiceImpl) PostFile(req *domain.ChatMessage,roomId string,ctx *gin.Context) (error) {

	fmt.Println("Upload File")

	fileType := req.File.Header.Get("Content-Type")

	file,err := req.File.Open()
	if err!= nil {
        return err
    }

	fileName := req.File.Filename

	now := time.Now()
	resp := make(chan *uploader.UploadResult)

	go func() {
		res, err := u.cloud.Upload.Upload(ctx,file,uploader.UploadParams{
			Type: api.Upload,
			ResourceType: "auto",
			DisplayName: fileName,
			FilenameOverride: fileName,
			PublicID: roomId+now.Format("2006_01_02_T15_04_05"),
		})
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
        "message": fileName,
		"url" : response.SecureURL,
        "createdAt": req.CreatedAt,
        "type": stringType,
	}
	_,_,err = colRef.Doc(roomId).Collection("chats").Add(ctx, content)
	if err!= nil {
        return err
    }
	return nil

}
