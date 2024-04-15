package service

import (
	"medichat-be/domain"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type ChatService interface {
	Post(req *domain.ChatMessage,roomId string,ctx *gin.Context) (error)}

type chatServiceImpl struct {
	client *firestore.Client

}

func NewChatServiceImpl(client  *firestore.Client) *chatServiceImpl {
	return &chatServiceImpl{
		client: client,
	}
}

func (u *chatServiceImpl) Post(req *domain.ChatMessage,roomId string,ctx *gin.Context) (error) {


	colRef := u.client.Collection("rooms");
	content := map[string]interface{}{
        "userId": req.UserId,
        "userName": req.UserName,
        "message": req.Message,
        "createdAt": req.CreatedAt,
        "type": req.Type,
	}
	_,_,err := colRef.Doc(roomId).Collection("chats").Add(ctx, content)
	if err!= nil {
        return err
    }
	return nil

}
