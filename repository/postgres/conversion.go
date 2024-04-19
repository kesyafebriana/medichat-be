package postgres

import "database/sql"

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
