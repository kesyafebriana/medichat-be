package postgis

import (
	"encoding/binary"
	"errors"
	"math"
)

var (
	ErrInvalidType = errors.New("postgis: invalid geometry type")
	ErrIncomplete  = errors.New("postgis: incomplete byte representation")
)

const (
	TypePoint = 0x00000001
	TypeMask  = 0x00FFFFFF
	FlagSRID  = 0x20000000
)

type EWKB struct {
	endian byte
	gType  uint32
	srid   uint32
	data   []byte
	coords []float64
}

func NewEWKB(b []byte, dataLen uint) (EWKB, error) {
	var ret EWKB
	l := len(b)

	if l < 9 {
		return EWKB{}, ErrIncomplete
	}

	ret.endian = b[0]

	var bo binary.ByteOrder

	if ret.endian == 0 {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}

	ret.gType = bo.Uint32(b[1:5])

	i := 5

	if ret.gType&FlagSRID != 0 {
		ret.srid = bo.Uint32(b[i : i+4])
		i += 4
	}

	ret.gType = ret.gType & TypeMask

	if i+int(dataLen) > l {
		return EWKB{}, ErrIncomplete
	}

	ret.data = b[i : i+int(dataLen)]
	i += int(dataLen)

	cl := (l - i) / 8
	ret.coords = make([]float64, cl)
	for j := 0; j < cl; j++ {
		ret.coords[j] = math.Float64frombits(bo.Uint64(b[i : i+8]))
		i += 8
	}

	return ret, nil
}
