package encode

import (
	"encoding/binary"
	"math"

	"github.com/et-nik/binngo/binn"
)

type intType uint8

func Uint(v uint) []byte {
	switch detectUintType(v) {
	case binn.Uint8Type:
		return Uint8(uint8(v))
	case binn.Uint16Type:
		return Uint16(uint16(v))
	case binn.Uint32Type:
		return Uint32(uint32(v))
	default:
		return Uint64(uint64(v))
	}
}

func Int(v int) []byte {
	t := detectIntType(v)

	result := []byte{}

	switch t {
	case binn.Int8Type:
		result = Int8(int8(v))
	case binn.Uint8Type:
		result = Uint8(uint8(v))
	case binn.Int16Type:
		result = Int16(int16(v))
	case binn.Uint16Type:
		result = Uint16(uint16(v))
	case binn.Int32Type:
		result = Int32(int32(v))
	case binn.Uint32Type:
		result = Uint32(uint32(v))
	case binn.Int64Type:
		result = Int64(int64(v))
	case binn.Uint64Type:
		result = Uint64(uint64(v))
	}

	return result
}

func detectUintType(v uint) intType {
	switch {
	case v <= math.MaxUint8:
		return binn.Uint8Type
	case v <= math.MaxUint16:
		return binn.Uint16Type
	case v <= math.MaxUint32:
		return binn.Uint32Type
	default:
		return binn.Uint64Type
	}
}

func detectIntType(v int) intType {
	t := binn.Int64Type

	if v > 0 {
		switch t {
		case binn.Int64Type:
			t = binn.Uint64Type
		case binn.Int32Type:
			t = binn.Uint32Type
		case binn.Int16Type:
			t = binn.Uint16Type
		case binn.Int8Type:
			t = binn.Uint8Type
		}
	}

	if t == binn.Int64Type ||
		t == binn.Int32Type ||
		t == binn.Int16Type {
		switch {
		case v >= math.MinInt8:
			t = binn.Int8Type
		case v >= math.MinInt16:
			t = binn.Int16Type
		case v >= math.MinInt32:
			t = binn.Int32Type
		}
	}

	if t == binn.Uint64Type ||
		t == binn.Uint32Type ||
		t == binn.Uint16Type {
		switch {
		case v <= math.MaxUint8:
			t = binn.Uint8Type
		case v <= math.MaxUint16:
			t = binn.Uint16Type
		case v <= math.MaxInt32:
			t = binn.Uint32Type
		}
	}

	return intType(t)
}

func Int8(v int8) []byte {
	return []byte{uint8(v)}
}

func Uint8(v uint8) []byte {
	return []byte{v}
}

func Uint16(v uint16) []byte {
	t := make([]byte, 2)
	binary.BigEndian.PutUint16(t, v)

	var r []byte
	r = append(r, t...)

	return r
}

func Int16(v int16) []byte {
	t := make([]byte, 2)
	binary.BigEndian.PutUint16(t, uint16(v))

	var r []byte
	r = append(r, t...)

	return r
}

func Uint32(v uint32) []byte {
	t := make([]byte, 4)
	binary.BigEndian.PutUint32(t, v)

	var r []byte
	r = append(r, t...)

	return r
}

func Int32(v int32) []byte {
	t := make([]byte, 4)
	binary.BigEndian.PutUint32(t, uint32(v))

	var r []byte
	r = append(r, t...)

	return r
}

func Uint64(v uint64) []byte {
	t := make([]byte, 8)
	binary.BigEndian.PutUint64(t, v)

	return t
}

func Int64(v int64) []byte {
	t := make([]byte, 8)
	binary.BigEndian.PutUint64(t, uint64(v))

	return t
}

func Size(size int, totalSize bool) []byte {
	sz := size

	if totalSize {
		sz++
	}

	if sz <= math.MaxInt8 {
		return []byte{byte(sz)}
	}

	if totalSize {
		sz += 3
	}

	return encodeSize32(sz)
}

func encodeSize32(s int) []byte {
	i := s | (-1 << 31)

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))

	return b
}
