package postgres

import (
	"context"
	"medichat-be/domain"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type chatRepository struct {
	querier Querier
}

func (r *chatRepository) GetChats(ctx context.Context, roomId int64) ([]domain.Chat, error) {
	sb := strings.Builder{}
	args := pgx.NamedArgs{}

	sb.WriteString(`
		SELECT ` + chatsColumns + `
		FROM chat_items c
		WHERE c.deleted_at IS NULL
	`)

	return queryFull(
		r.querier, ctx, sb.String(),
		scanChats,
		args,
	)
}

func (r *chatRepository) AddChat(ctx context.Context, chat domain.Chat) (domain.Chat, error) {
	q := `
		INSERT INTO chat_items(`+chatsColumns+`)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		`
	return queryOneFull(
		r.querier, ctx, q,
		scanChats,
		chat.ID, chat.RoomId, chat.Type, chat.Message, chat.File, chat.UserId, chat.UserName,
	)
}

func (r *chatRepository) AddRoom(ctx context.Context, UserId int,DoctorId int,EndAt time.Time ) (domain.Room, error) {
	q := `
		INSERT INTO chat_rooms(`+roomsColumns+`)
		VALUES
		($1, $2, $3)
		RETURNING id, user_id, doctor_id, end_at
		`
	return queryOneFull(
		r.querier, ctx, q,
		scanRooms,
		UserId, DoctorId, EndAt,
	)
}