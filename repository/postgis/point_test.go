package postgis_test

import (
	"encoding/hex"
	"fmt"
	"medichat-be/repository/postgis"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPointFromEWKB(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    postgis.Point
		wantErr error
	}{
		{
			name: "should correctly parse EWKB",
			hex:  "0101000020E6100000CDCCCCCCCC8C5C40295C8FC2F5A824C0",
			want: postgis.NewPoint(114.20, -10.33),
		},
		{
			name: "should correctly parse EWKB big endian",
			hex:  "0020000001000010E6405C8CCCCCCCCCCDC024A8F5C28F5C29",
			want: postgis.NewPoint(114.20, -10.33),
		},
		{
			name:    "should return error when parsing with incorrect type",
			hex:     "0102000020E6100000CDCCCCCCCC8C5C40295C8FC2F5A824C0",
			wantErr: postgis.ErrInvalidType,
		},
		{
			name:    "should return error when parsing incomplete data",
			hex:     "0101000020E6100000CDCCCCCCCC8C5C40",
			wantErr: postgis.ErrIncomplete,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			b, _ := hex.DecodeString(tt.hex)

			// when
			got, err := postgis.NewPointFromEWKB(b)

			// then
			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestPoint_Scan(t *testing.T) {
	tests := []struct {
		name    string
		hex     any
		want    postgis.Point
		wantErr error
	}{
		{
			name: "should correctly parse EWKB",
			hex:  "0101000020E6100000CDCCCCCCCC8C5C40295C8FC2F5A824C0",
			want: postgis.NewPoint(114.20, -10.33),
		},
		{
			name:    "should return error when parsing incomplete data",
			hex:     "0101000020E6100000CDCCCCCCCC8C5C40",
			wantErr: postgis.ErrIncomplete,
		},
		{
			name:    "should return error when scanning an integer",
			hex:     1002,
			wantErr: postgis.ErrInvalidType,
		},
		{
			name:    "should return error when parsing with invalid hex",
			hex:     "!30a158",
			wantErr: fmt.Errorf("encoding/hex: invalid byte: U+0021 '!'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			var p postgis.Point

			// when
			err := p.Scan(tt.hex)

			// then
			assert.Equal(t, tt.want, p)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			assert.Nil(t, err)
		})
	}
}
