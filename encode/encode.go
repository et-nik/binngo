package encode

import (
	"encoding"
	"reflect"
	"sync"

	"github.com/et-nik/binngo/binn"
)

type encoderFunc func(v reflect.Value) ([]byte, error)

var encoderCache sync.Map // map[reflect.Type]encoderFunc

var (
	marshalerType     = reflect.TypeOf((*Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

type Marshaler interface {
	MarshalBINN() ([]byte, error)
}

func marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)

	if !rv.IsValid() {
		return nil, ErrInvalidValue
	}

	if !rv.IsValid() {
		return nil, ErrInvalidValue
	}

	enc := loadEncodeFunc(rv.Type())

	return enc(rv)
}

func loadEncodeFunc(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(encoderFunc)
	}

	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	fi, loaded := encoderCache.LoadOrStore(t, encoderFunc(func(v reflect.Value) ([]byte, error) {
		wg.Wait()
		return f(v)
	}))
	if loaded {
		return fi.(encoderFunc)
	}

	f = newTypeEncoder(t)
	wg.Done()
	encoderCache.Store(t, f)
	return f
}

func newTypeEncoder(t reflect.Type) encoderFunc {
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}
	if t.Implements(textMarshalerType) {
		return textMarshalerEncoder
	}

	switch t.Kind() {
	case reflect.Bool:
		return func(v reflect.Value) ([]byte, error) {
			if v.Bool() {
				return []byte{binn.True}, nil
			}
			return []byte{binn.False}, nil
		}
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Interface:
		return func(v reflect.Value) ([]byte, error) {
			if v.IsNil() {
				return []byte{binn.Null}, nil
			}

			return loadEncodeFunc(v.Elem().Type())(v.Elem())
		}
	case reflect.String:
		return func(v reflect.Value) ([]byte, error) {
			var bytes []byte

			bytes = append(bytes, Uint8(binn.StringType)...)
			bytes = append(bytes, String(v.String())...)
			bytes = append(bytes, 0x00)

			return bytes, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(v reflect.Value) ([]byte, error) {
			var bytes []byte

			bytes = append(bytes, Uint8(uint8(detectIntType(int(v.Int()))))...)
			bytes = append(bytes, Int(int(v.Int()))...)

			return bytes, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return func(v reflect.Value) ([]byte, error) {
			var bytes []byte

			bytes = append(bytes, Uint8(uint8(detectUintType(uint(v.Uint()))))...)
			bytes = append(bytes, Uint(uint(v.Uint()))...)

			return bytes, nil
		}
	case reflect.Float32:
		return func(v reflect.Value) ([]byte, error) {
			var bytes []byte

			bytes = append(bytes, Uint8(binn.Float32Type)...)
			bytes = append(bytes, Float32(float32(v.Float()))...)

			return bytes, nil
		}
	case reflect.Float64:
		return func(v reflect.Value) ([]byte, error) {
			var bytes []byte

			bytes = append(bytes, Uint8(binn.Float64Type)...)
			bytes = append(bytes, Float64(v.Float())...)

			return bytes, nil
		}
	case reflect.Slice, reflect.Array:
		return newArrayEncoder(t)
	case reflect.Ptr:
		return newPtrEncoder(t)
	}

	return func(v reflect.Value) ([]byte, error) {
		return nil, &UnsupportedTypeError{t}
	}
}

func marshalerEncoder(v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return []byte{0x00}, nil
	}
	m, ok := v.Interface().(Marshaler)
	if !ok {
		return []byte{0x00}, nil
	}
	b, err := m.MarshalBINN()
	if err != nil {
		return nil, &MarshalerError{v.Type(), err, "MarshalBINN"}
	}

	return b, nil
}
