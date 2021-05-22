package encode

import "reflect"

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "binn: unsupported type: " + e.Type.String()
}
