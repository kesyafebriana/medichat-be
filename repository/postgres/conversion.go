package postgres

import (
	"database/sql"
)

func fromStringPtr(s *string) sql.NullString {
	var ret sql.NullString
	if s != nil {
		ret.Valid, ret.String = true, *s
	}
	return ret
}

func toStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func fromInt64Ptr(i *int64) sql.NullInt64 {
	var ret sql.NullInt64
	if i != nil {
		ret.Valid, ret.Int64 = true, *i
	}
	return ret
}

func toInt64Ptr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}
