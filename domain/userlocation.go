package domain

type UserLocation struct {
	ID   int64
	User User

	Alias    string
	Name     string
	PhotoURL string
	IsActive bool
}
