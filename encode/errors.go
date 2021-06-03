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

type MarshalerError struct {
	Type       reflect.Type
	Err        error
	sourceFunc string
}

func (e *MarshalerError) Error() string {
	srcFunc := e.sourceFunc
	if srcFunc == "" {
		srcFunc = "MarshalBINN"
	}
	return "binn: error calling " + srcFunc +
		" for type " + e.Type.String() +
		": " + e.Err.Error()
}

func (e *MarshalerError) Unwrap() error { return e.Err }
