// Package decode implements BINN decoding.
package decode

import (
	"bytes"
	"reflect"
)

type Unmarshaler interface {
	UnmarshalBINN([]byte) error
}

var unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

func Unmarshal(data []byte, v interface{}) error {
	decoder := NewDecoder(bytes.NewReader(data))

	rt := reflect.TypeOf(v)

	if rt.Implements(unmarshalerType) {
		err := decodeUnmarshaler(data, v)
		if err != nil {
			return err
		}

		return nil
	}

	return decoder.Decode(v)
}

func decodeUnmarshaler(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)

	unm, ok := rv.Interface().(Unmarshaler)
	if !ok {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	err := unm.UnmarshalBINN(data)
	if err != nil {
		return err
	}

	return nil
}
