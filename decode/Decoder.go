package decode

import (
	"bytes"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	r := bytes.NewReader(data)
	return decode(r, v)
}
