package postgis

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"medichat-be/domain"
)

type Point struct {
	X float64
	Y float64
}

func NewPoint(x, y float64) Point {
	return Point{
		X: x,
		Y: y,
	}
}

func NewPointFromCoordinate(c domain.Coordinate) Point {
	return Point{
		X: c.Longitude,
		Y: c.Latitude,
	}
}

func (p Point) ToCoordinate() domain.Coordinate {
	return domain.Coordinate{
		Longitude: p.X,
		Latitude:  p.Y,
	}
}

func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%f %f)", p.X, p.Y), nil
}

func (p *Point) Scan(value any) error {
	switch v := value.(type) {
	case string:
		b, err := hex.DecodeString(v)
		if err != nil {
			return err
		}
		po, err := NewPointFromEWKB(b)
		if err != nil {
			return err
		}
		*p = po
		return nil
	default:
		return ErrInvalidType
	}
}

func NewPointFromEWKB(b []byte) (Point, error) {
	ewkb, err := NewEWKB(b, 0)
	if err != nil {
		return Point{}, err
	}

	if ewkb.gType != TypePoint {
		return Point{}, ErrInvalidType
	}

	if len(ewkb.coords) < 2 {
		return Point{}, ErrIncomplete
	}

	return Point{
		X: ewkb.coords[0],
		Y: ewkb.coords[1],
	}, nil
}
