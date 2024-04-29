package domain

type Rating struct {
	ID     int64
	User   User
	Doctor Doctor

	Name    string
	IsLiked bool
}
