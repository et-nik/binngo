package decode

import (
	"encoding/binary"

	"github.com/et-nik/binngo/binn"
)

func DecodeType(b []byte) binn.BinnType {
	return binn.BinnType(DecodeUint8(b))
}

func DecodeUint8(b []byte) uint8 {
	return b[0]
}

func DecodeUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func DecodeUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func DecodeUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func DecodeInt8(b []byte) int8 {
	return int8(b[0])
}

func DecodeInt16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b))
}

func DecodeInt32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

func DecodeInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func DecodeString(b []byte) string {
	return string(b)
}
