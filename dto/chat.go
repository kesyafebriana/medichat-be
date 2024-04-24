package dto

import (
	"mime/multipart"
	"time"
)

type ChatRoom struct{
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
	DoctorId int `json:"doctorId"`
	DoctorName string `json:"doctorName"`
	IsTyping []string `json:"isTyping"`
	UserId int `json:"userId"`
	UserName string `json:"userName"`
}

type ChatMessage struct{
	CreatedAt time.Time `json:"createdAt"`
	Message string `json:"message"`
	Type string `json:"type"`
	UserId int `json:"userId"`
	UserName string `json:"userName"`
	File *multipart.FileHeader `json:"file"`
}

