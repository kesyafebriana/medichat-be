package domain

import "time"

type Doctor struct {
	ID             int64
	Specialization Specialization

	STR           string
	WorkLocation  string
	Gender        string
	PhoneNumber   string
	IsActive      bool
	StartWorkDate time.Time
	Price         int64
}
