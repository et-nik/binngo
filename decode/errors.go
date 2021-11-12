package decode

import (
	"errors"
	"reflect"
	"strings"
)

var (
	ErrUnknownType        = errors.New("unknown storage type")
	ErrCantSetValue       = errors.New("can't set value")
	ErrItemNotFound       = errors.New("item not found")
	ErrInvalidItem        = errors.New("invalid item")
	ErrInvalidStructValue = errors.New("invalid struct value")
)

type FailedToReadTypeError struct {
	Previous error
}

func (err *FailedToReadTypeError) Error() string {
	text := strings.Builder{}
	text.WriteString("failed to read type")

	if err.Previous == nil {
		return text.String()
	}

	text.WriteString(": ")
	text.WriteString(err.Previous.Error())

	return text.String()
}

func (err *FailedToReadTypeError) Unwrap() error {
	return err.Previous
}

type FailedToReadSizeError struct {
	Previous error
}

func (err *FailedToReadSizeError) Error() string {
	text := strings.Builder{}
	text.WriteString("failed to read size")

	if err.Previous == nil {
		return text.String()
	}

	text.WriteString(": ")
	text.WriteString(err.Previous.Error())

	return text.String()
}

func (err *FailedToReadSizeError) Unwrap() error {
	return err.Previous
}

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
