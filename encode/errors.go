package encode

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidValue = errors.New("invalid value")
)

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "binn: unsupported type: " + e.Type.String()
}
