package encode

import (
	"encoding/binary"
	"math"
)

func Float32(f float32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(f))

	return b
}

func Float64(f float64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))

	return b
}
