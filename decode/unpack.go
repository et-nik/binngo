package decode

import (
	"encoding/binary"

	"github.com/et-nik/binngo/binn"
)

func Type(b []byte) binn.Type {
	return binn.Type(Uint8(b))
}

func Uint8(b []byte) uint8 {
	return b[0]
}

func Uint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func Uint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func Uint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func Int8(b []byte) int8 {
	return int8(b[0])
}

func Int16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b))
}

func Int32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

func Int64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func String(b []byte) string {
	return string(b)
}
