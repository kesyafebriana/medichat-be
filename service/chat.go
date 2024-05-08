package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"medichat-be/constants"
	"medichat-be/domain"
	"medichat-be/dto"
	"medichat-be/util"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/pdfcrowd/pdfcrowd-go"
)

type ChatService interface {
	PostMessage(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error)
	PostFile(req *dto.ChatMessage,roomId string,ctx *gin.Context) (error)
	CreateRoom(doctorId int,ctx *gin.Context) (error)
	CloseRoom(roomId string,ctx *gin.Context) (error)
	CreateNote(userId int64, roomId,message string,ctx *gin.Context) (error)
	Prescribe(req *dto.ChatPrescription,roomId string,ctx *gin.Context) (error)
}

type chatService struct {
	dataRepository domain.DataRepository
	client  *firestore.Client
	cloud util.CloudinaryProvider
	doctorNoteTemplate *template.Template
}

type ChatServiceOpts struct {
	DataRepository domain.DataRepository
	Client  *firestore.Client
	Cloud util.CloudinaryProvider
}
func NewChatService(opts ChatServiceOpts) *chatService {


	doctorNoteTemplate, err := template.ParseFiles("templates/doctor-notes.html")
	if err != nil {
		return nil
	}

	return &chatService{
		dataRepository: opts.DataRepository,
		doctorNoteTemplate: doctorNoteTemplate,
		client: opts.Client,
		cloud: opts.Cloud,
	}
}

func (u *chatService) Prescribe(req *dto.ChatPrescription,roomId string,ctx *gin.Context) (error) {

	userRepository := u.dataRepository.UserRepository()
	doctorRepository := u.dataRepository.DoctorRepository()
	productRepository := u.dataRepository.ProductRepository()

	doctorId,err := util.GetAccountIDFromContext(ctx);
	if err!= nil {
		return err
    }
	doctor,err:= doctorRepository.GetByAccountID(ctx,doctorId)
	if err!= nil {
		return err
    }

	user,err := userRepository.GetByID(ctx,int64(req.UserId))
	if err!= nil {
        return err
    }
	now := time.Now()

	var drugs []map[string]interface{}

	for i := 0; i < len(req.Drugs); i++ {
		prod,err:= productRepository.GetById(ctx,int64(req.Drugs[i].ProductId))
		if err!= nil {
            return err
        }
		drugs = append(drugs, map[string]interface{}{
			"id": req.Drugs[i].ProductId,
			"name": prod.Name,
            "count": req.Drugs[i].Count,
            "direction": req.Drugs[i].Direction,
            "picture": prod.Picture,
		})
	}


	prescription := map[string]interface{}{
		"userId":req.UserId,
		"drugs":drugs,
	}

	json,err := json.Marshal(prescription)
	if err!= nil {
        return err
    }

	colRef := u.client.Collection("rooms");
	content := map[string]interface{}{
        "userId": doctor.Account.ID,
        "userName": user.Account.Name,
        "message": string(json),
        "createdAt": now,
		"url" : "",
        "type": "message/prescription",
	}
	_,_,err = colRef.Doc(roomId).Collection("chats").Add(ctx, content)
	if err!= nil {
        return err
    }
	return nil
}

func (u *chatService) CreateNote(userId int64, roomId,message string,ctx *gin.Context) (error){

	doctorRepository := u.dataRepository.DoctorRepository()
	userRepository := u.dataRepository.UserRepository()

	user,err := userRepository.GetByID(ctx,userId)

	if (err != nil){
        return err
    }

	doctorId, err := util.GetAccountIDFromContext(ctx)
	if err != nil {
		return err

	}

	doctor,err := doctorRepository.GetByAccountID(ctx,doctorId)
	if (err != nil){
		return err
	}

	nowString := time.Now().Format("02-01-2006 : 03:04:05")

	name := user.Account.Name;
	dob := user.DateOfBirth.Format("02-01-2006")

	var body bytes.Buffer
	u.doctorNoteTemplate.Execute(&body, struct {
		Fullname   		string
		DateOfBirth     string
		Date 			string
		DoctorName 		string
		DoctorNumber 	string
		DoctorMessage 	string
	}{
		Fullname:   name,
		DateOfBirth: dob,
		Date: nowString,
		DoctorMessage: message,
		DoctorName: doctor.Account.Name,
		DoctorNumber: doctor.STR,
	})

	client := pdfcrowd.NewHtmlToPdfClient("demo", "ce544b6ea52a5621fb9d55f8b542d14d")

	var pdf bytes.Buffer

    err = client.ConvertStringToStream(body.String(),&pdf)
	if err!= nil {
        return err
    }

	fileName := "Doctor."+doctor.Account.Name+"Notes"+user.Account.Name+nowString

	file := bytes.NewReader(pdf.Bytes())

	opts := util.SendFileOpts{
		Context: ctx,
		Filename: fileName,
		Roomid: roomId,
		File: file,
	}
	response, err := u.cloud.SendFile(opts)

	if err!= nil {
        return err
    }

	colRef := u.client.Collection("rooms");

	content := map[string]interface{}{
		"userId": doctor.Account.ID,
		"userName": doctor.Account.Name,
		"message": fileName,
		"url" : response.SecureURL,
		"createdAt": time.Now(),
		"type": "message/pdf",
	}
	_,_,err = colRef.Doc(roomId).Collection("chats").Add(ctx, content)
	if err!= nil {
		return err
	}

	if err!= nil {
        return err
    }

	return nil
}


func (u *chatService) CreateRoom(doctorId int,ctx *gin.Context) (error) {

	var req dto.ChatRoom

	req.DoctorId = doctorId

	chatRepository := u.dataRepository.ChatRepository()
	userRepository :=  u.dataRepository.UserRepository()
	doctorRepository := u.dataRepository.DoctorRepository()


	userId, err := util.GetAccountIDFromContext(ctx)
	if err != nil {
        return err
	}
	user,err:= userRepository.GetByAccountID(ctx,userId)

	if err!= nil {
        return err
    }

	req.UserId = int(user.ID)
	req.UserName = user.Account.Name

	doctor,err := doctorRepository.GetByID(ctx,int64(req.DoctorId))
	if err != nil {
		return err
	}
	req.DoctorName = doctor.Account.Name


	date:= time.Now()
	req.Start = date

	extra , _ := time.ParseDuration(constants.ChatDuration)

	req.End = date.Add(extra)

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

	tar := map[string]interface{}{
		"doctorId":req.DoctorId,
		"doctorName":req.DoctorName,
		"end": req.End,
		"start": req.Start,
		"userId":req.UserId,
		"userName":req.UserName,
		"open" : true,
		"isTyping":req.IsTyping,
	}

	_,err = colRef.Doc(roomId).Set(ctx,tar)

	if err!= nil {
        return err
    }

	return nil

}

func (u *chatService) CloseRoom(roomId string,ctx *gin.Context) (error) {

	chatRepository := u.dataRepository.ChatRepository()

	colRef := u.client.Collection("rooms");

	_,err := colRef.Doc(roomId).Update(ctx,[]firestore.Update{
		{Path: "end", Value: time.Now()},
	})
	if err!= nil {
        return err
    }

	docs,err := colRef.Doc(roomId).Collection("chats").Documents(ctx).GetAll()
	if err !=nil{
		return err
	}

	room_id,err := strconv.Atoi(roomId)
	if (err!=nil){
		return err
	}

	for i := 0; i < len(docs); i++ {
		data := docs[i].Data()


		userId:= int(data["userId"].(int64))

		if err !=nil {
			return err
		}

		fmt.Println(data)

		chat:= domain.Chat{
			RoomId: int64(room_id),
            Message: data["message"].(string),
            File: data["url"].(string),
            Type: data["type"].(string),
            UserId: userId,
			UserName: data["userName"].(string),
			CreatedAt: data["createdAt"].(time.Time),
		}

		_,err := chatRepository.AddChat(ctx, chat)
		if err!= nil {
            return err
        }
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
