package domain

import (
	"time"
)

type ChatRoom struct{
	Date time.Time `json:"date"`
	DoctorId int `json:"doctorId"`
	DoctorName string `json:"doctorName"`
	IsTyping []string `json:"isTyping"`
	Open bool `json:"open"`
	UserId int `json:"userId"`
	UserName string `json:"userName"`
}

type ChatMessage struct{
	CreatedAt time.Time `json:"createdAt"`
	Message string `json:"message"`
	Type string `json:"type"`
	UserId int `json:"userId"`
	UserName string `json:"userName"`
}

