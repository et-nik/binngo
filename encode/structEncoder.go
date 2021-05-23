package encode

import (
	"reflect"
	"strings"

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
		keyName := strings.Split(se.t.Field(i).Tag.Get("binn"), ",")[0]
		if keyName == "" {
			keyName = se.t.Field(i).Name
		}

		dataBytes = append(dataBytes, String(keyName)...)

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
