package encode

import (
	"encoding"
	"reflect"

	"github.com/et-nik/binngo/binn"
)

func newMapEncoder(t reflect.Type) encoderFunc {
	switch t.Key().Kind() {
	case reflect.String:
		me := mapObjectEncoder{
			newTypeEncoder(t.Elem()),
			t,
		}
		return me.encode
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		me := mapEncoder{newTypeEncoder(t.Elem())}
		return me.encode
	default:
		if t.Key().Implements(textMarshalerType) {
			me := mapObjectEncoder{
				newTypeEncoder(t.Elem()),
				t,
			}
			return me.encode
		}

		return func(v reflect.Value) ([]byte, error) {
			return nil, &UnsupportedTypeError{t}
		}
	}
}

type mapObjectEncoder struct {
	elemEnc  encoderFunc
	itemType reflect.Type
}

func (me *mapObjectEncoder) encode(v reflect.Value) ([]byte, error) {
	dataBytes := []byte{}

	keys := v.MapKeys()

	for _, key := range keys {
		bval, err := encodeTextKey(key)
		if err != nil {
			return nil, err
		}
		dataBytes = append(dataBytes, bval...)

		encodedValue, err := me.elemEnc(v.MapIndex(key))
		if err != nil {
			return nil, err
		}

		dataBytes = append(dataBytes, encodedValue...)
	}

	bytes := []byte{}

	typeBytes := Uint8(binn.ObjectType)
	countBytes := Size(len(keys), false)

	bytes = append(bytes, typeBytes...)
	bytes = append(bytes, Size(len(typeBytes)+len(dataBytes)+len(countBytes), true)...)
	bytes = append(bytes, countBytes...)
	bytes = append(bytes, dataBytes...)

	return bytes, nil
}

type mapEncoder struct {
	elemEnc encoderFunc
}

func (me *mapEncoder) encode(v reflect.Value) ([]byte, error) {
	var dataBytes []byte

	iter := v.MapRange()

	for iter.Next() {
		key := iter.Key()
		val := iter.Value()

		dataBytes = append(dataBytes, Int32(int32(key.Int()))...)

		encodedValue, err := me.elemEnc(val)
		if err != nil {
			return nil, err
		}
		dataBytes = append(dataBytes, encodedValue...)
	}

	keys := v.MapKeys()

	bytes := []byte{}

	typeBytes := Uint8(binn.MapType)
	countBytes := Size(len(keys), false)

	bytes = append(bytes, typeBytes...)
	bytes = append(bytes, Size(len(typeBytes)+len(dataBytes)+len(countBytes), true)...)
	bytes = append(bytes, countBytes...)
	bytes = append(bytes, dataBytes...)

	return bytes, nil
}

func encodeTextKey(v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.String {
		return String(v.String()), nil
	}

	m, ok := v.Interface().(encoding.TextMarshaler)
	if !ok {
		return []byte{binn.Null}, nil
	}

	s, err := m.MarshalText()
	if err != nil {
		return nil, err
	}

	return String(string(s)), nil
}
