package binn

import (
	"github.com/et-nik/binngo/decode"
	"github.com/et-nik/binngo/encode"
)

func Marshal(v interface{}) ([]byte, error) {
	return encode.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return decode.Unmarshal(data, v)
}
