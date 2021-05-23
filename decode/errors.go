package decode

import (
	"errors"
	"reflect"
)

var (
	ErrUnknownType        = errors.New("unknown storage type")
	ErrCantSetValue       = errors.New("can't set value")
	ErrItemNotFound       = errors.New("item not found")
	ErrInvalidItem        = errors.New("invalid item")
	ErrInvalidStructValue = errors.New("invalid struct value")
	ErrFailedToReadType   = errors.New("failed to read type")
	ErrFailedToReadSize   = errors.New("failed to read size")
)

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "binn: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "binn: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "binn: Unmarshal(nil " + e.Type.String() + ")"
}

type UnknownValueError struct {
	Expected reflect.Kind
	Got      reflect.Kind
}

func (e *UnknownValueError) Error() string {
	return "binn: Unknown value. Expected " + e.Expected.String() + ", got " + e.Got.String()
}
