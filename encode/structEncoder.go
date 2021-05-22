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
			EncodeString(se.t.Field(i).Name)...
		)

		encodeValue := loadEncodeFunc(v.Field(i).Type())
		val, err := encodeValue(v.Field(i))
		if err != nil {
			return nil, err
		}

		dataBytes = append(dataBytes, val...)
	}

	bytes := []byte{}

	typeBytes := EncodeUint8(binn.ObjectType)
	countBytes := EncodeSize(v.NumField(), false)

	bytes = append(bytes, typeBytes...)
	bytes = append(bytes, EncodeSize(len(typeBytes) + len(dataBytes) + len(countBytes), true)...)
	bytes = append(bytes, countBytes...)
	bytes = append(bytes, dataBytes...)

	return bytes, nil
}
