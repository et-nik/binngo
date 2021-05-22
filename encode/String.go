package encode

import (
	"encoding"
	"reflect"

	"github.com/et-nik/binngo/binn"
)

func EncodeString(s string) []byte {
	var t []byte

	t = append(t, EncodeSize(len(s), false)...)
	t = append(t, []byte(s)...)

	return t
}

func textMarshalerEncoder(v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return []byte{binn.Null}, nil
	}

	m, ok := v.Interface().(encoding.TextMarshaler)
	if !ok {
		return []byte{binn.Null}, nil
	}
	b, err := m.MarshalText()
	if err != nil {
		return nil, err
	}

	return EncodeString(string(b)), nil
}
