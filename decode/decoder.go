// Package decode implements BINN decoding.
package decode

import (
	"bytes"
)

func Unmarshal(data []byte, v interface{}) error {
	decoder := NewDecoder(bytes.NewReader(data))

	return decoder.Decode(v)
}
