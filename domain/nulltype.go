package domain

import "time"

type NullInt struct {
	Int   int
	Valid bool
}

func NewNullInt(i int) NullInt {
	return NullInt{
		Int:   i,
		Valid: true,
	}
}

func FromIntPtr(i *int) NullInt {
	if i == nil {
		return NullInt{}
	}
	return NullInt{
		Int:   *i,
		Valid: true,
	}
}

func (i *NullInt) ToIntPtr() *int {
	if !i.Valid {
		return nil
	}
	return &i.Int
}

type NullString struct {
	String string
	Valid  bool
}

func NewNullString(s string) NullString {
	return NullString{
		String: s,
		Valid:  true,
	}
}

func FromStringPtr(s *string) NullString {
	if s == nil {
		return NullString{}
	}
	return NullString{
		String: *s,
		Valid:  true,
	}
}

func (s *NullString) ToStringPtr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{
		Time:  t,
		Valid: true,
	}
}

func FromTimePtr(t *time.Time) NullTime {
	if t == nil {
		return NullTime{}
	}
	return NullTime{
		Time:  *t,
		Valid: true,
	}
}

func (t *NullTime) ToTimePtr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
