package domain

import (
	"context"
	"time"
)

type Chat struct {
	RoomId 			int64
	Message 		string
	File 			string
	Type 			string
	UserId 			int
	UserName 		string
}

type Room struct{
	ID 			int64
	UserId 		int64
	DoctorId 	int64
	EndAt  		time.Time
}


type ChatRepository interface {
	GetChats(ctx context.Context, roomId int64) ([]Chat, error)
	AddChat(ctx context.Context, chat Chat) (Chat, error)
	AddRoom(ctx context.Context, UserId int,DoctorId int,EndAt time.Time ) (Room, error)
}

