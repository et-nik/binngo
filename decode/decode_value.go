package decode

import (
	"io"
	"reflect"

	"github.com/et-nik/binngo/binn"
)

type valueDecoder struct {
	binnType binn.Type
}

func newValueDecoder(binnType binn.Type) *valueDecoder {
	return &valueDecoder{binnType}
}

func (vd *valueDecoder) DecodeValue(reader io.Reader, v interface{}) error {
	valuePtr := reflect.ValueOf(v)
	value := valuePtr.Elem()

	if !value.CanSet() {
		return ErrCantSetValue
	}

	bval, err := readValue(vd.binnType, reader)
	if err != nil {
		return err
	}

	converted, err := decodeItem(value.Type(), vd.binnType, bval)
	if err != nil {
		return err
	}

	if value.Kind() != reflect.ValueOf(converted).Kind() && value.Kind() != reflect.Interface {
		return &UnknownValueError{reflect.ValueOf(converted).Kind(), value.Kind()}
	}

	value.Set(reflect.ValueOf(converted))

	return nil
}
