package encode

import (
	"reflect"

	"github.com/et-nik/binngo/binn"
)

type structEncoder struct {
	t reflect.Type
}

func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{t}
	return se.encode
}

func (se *structEncoder) encode(v reflect.Value) ([]byte, error) {
	dataBytes := []byte{}

	for i := 0; i < v.NumField(); i++ {
		dataBytes = append(dataBytes,
			String(se.t.Field(i).Name)...,
		)

		encodeValue := loadEncodeFunc(v.Field(i).Type())
		val, err := encodeValue(v.Field(i))
		if err != nil {
			return nil, err
		}

		dataBytes = append(dataBytes, val...)
	}

	bytes := []byte{}

	typeBytes := Uint8(binn.ObjectType)
	countBytes := Size(v.NumField(), false)

	bytes = append(bytes, typeBytes...)
	bytes = append(bytes, Size(len(typeBytes)+len(dataBytes)+len(countBytes), true)...)
	bytes = append(bytes, countBytes...)
	bytes = append(bytes, dataBytes...)

	return bytes, nil
}
