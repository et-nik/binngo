package encode

import (
	"reflect"

	"github.com/et-nik/binngo/binn"
)

type arrayEncoder struct {
	elemEnc encoderFunc
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	ae := arrayEncoder{newTypeEncoder(t.Elem())}
	return ae.encode
}

func (ae *arrayEncoder) encode(v reflect.Value) ([]byte, error) {
	var dataBytes []byte
	n := v.Len()
	for i := 0; i < n; i++ {
		encoded, err := ae.elemEnc(v.Index(i))
		if err != nil {
			return nil, err
		}

		dataBytes = append(dataBytes, encoded...)
	}

	bytes := []byte{}

	typeBytes := EncodeUint8(binn.ListType)
	countBytes := EncodeSize(v.Len(), false)

	bytes = append(bytes, typeBytes...)
	bytes = append(bytes, EncodeSize(len(typeBytes) + len(dataBytes) + len(countBytes), true)...)
	bytes = append(bytes, countBytes...)
	bytes = append(bytes, dataBytes...)

	return bytes, nil
}
