package encode

import (
	"reflect"

	"github.com/et-nik/binngo/binn"
)

type ptrEncoder struct {
	elemEnc encoderFunc
}

func newPtrEncoder(t reflect.Type) encoderFunc {
	enc := ptrEncoder{newTypeEncoder(t.Elem())}
	return enc.encode
}

func (pe ptrEncoder) encode(v reflect.Value) ([]byte, error) {
	if v.IsNil() {
		return []byte{binn.Null}, nil
	}

	return pe.elemEnc(v.Elem())
}
