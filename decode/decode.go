package decode

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/et-nik/binngo/binn"

	"github.com/cstockton/go-conv"
)

const (
	maxOneByteSize = 127
)

type decodeFunc func(reader io.Reader, v interface{}) error

var decoderCache sync.Map // map[binn.BinnType]decodeFunc

type readLen int

func decode(reader io.Reader, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	containerType, _, err := readType(reader)
	if err != nil {
		return err
	}

	err = decodeStorage(containerType, reader, v)
	if err != nil {
		return err
	}

	return nil
}

func decodeStorage(containerType binn.BinnType, reader io.Reader, v interface{}) error {
	decoder := loadDecodeFunc(containerType)
	return decoder(reader, v)
}

func decodeItem(rt reflect.Type, btype binn.BinnType, bval []byte) (interface{}, error) {
	var v interface{}
	var err error

	switch btype {
	case binn.Null:
		return nil, nil
	case binn.True:
		return true, nil
	case binn.False:
		return false, nil
	case binn.Uint8Type:
		v = DecodeUint8(bval)
	case binn.Uint16Type:
		v = DecodeUint16(bval)
	case binn.Uint32Type:
		v = DecodeUint32(bval)
	case binn.Uint64Type:
		v = DecodeUint64(bval)
	case binn.Int8Type:
		v = DecodeInt8(bval)
	case binn.Int16Type:
		v = DecodeInt16(bval)
	case binn.Int32Type:
		v = DecodeInt32(bval)
	case binn.Int64Type:
		v = DecodeInt64(bval)
	case binn.StringType:
		v = DecodeString(bval[:len(bval)-1])
	case binn.BlobType:
		v = bval
	case binn.ListType:
		var l []interface{}
		br := bytes.NewReader(bval)
		_, wasReadLen, _ := readSize(br)
		cnt, wasReadCnt, _ := readSize(br)
		wasRead := wasReadLen + wasReadCnt

		//et := reflect.TypeOf(rt)
		//l := reflect.MakeSlice(reflect.SliceOf(et), 0, cnt)

		err = decodeListItems(br, &l, len(bval)-int(wasRead), wasRead, cnt)
		if err != nil {
			return nil, err
		}

		return l, nil
	case binn.MapType:
		m := map[int]interface{}{}
		br := bytes.NewReader(bval)
		cnt, rl, _ := readSize(br)
		err = decodeMapItems(br, &m, len(bval)-int(rl), rl, cnt)
		if err != nil {
			return nil, err
		}

		return m, nil
	case binn.ObjectType:
		var obj interface{}
		if rt.Kind() == reflect.Interface {
			obj = map[string]interface{}{}
		} else {
			ptr := reflect.New(rt)
			obj = ptr.Interface()
		}

		br := bytes.NewReader(bval)

		err = decodeStorage(btype, br, &obj)

		if err != nil {
			return nil, err
		}

		return obj, nil
	}

	if rt.Kind() == reflect.Interface {
		v, err = convertToKind(kindMapper[btype], v)
		if err != nil {
			return nil, err
		}
	} else {
		v, err = convertToType(rt, v)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func loadDecodeFunc(bt binn.BinnType) decodeFunc {
	if fi, ok := decoderCache.Load(bt); ok {
		return fi.(decodeFunc)
	}

	var (
		wg sync.WaitGroup
		f decodeFunc
	)
	wg.Add(1)
	fi, loaded := decoderCache.LoadOrStore(bt, decodeFunc(func(reader io.Reader, v interface{}) error {
		wg.Wait()
		return f(reader, v)
	}))
	if loaded {
		return fi.(decodeFunc)
	}

	f = newTypeDecoder(bt)
	wg.Done()
	decoderCache.Store(bt, f)
	return f
}

func newTypeDecoder(bt binn.BinnType) decodeFunc {
	switch bt {
	case binn.ListType:
		return decodeList
	case binn.MapType:
		return decodeMap
	case binn.ObjectType:
		return decodeObject
	case binn.True, binn.False:
		return func(reader io.Reader, v interface{}) error {
			valuePtr := reflect.ValueOf(v)
			value := valuePtr.Elem()

			if value.Kind() != reflect.Bool && value.Kind() != reflect.Interface {
				return &UnknownValueError{reflect.Bool, value.Kind()}
			}

			if !value.CanSet() {
				return ErrCantSetValue
			}

			value.Set(reflect.ValueOf(bt == binn.True))

			return nil
		}
	case binn.Null:
		return func(reader io.Reader, v interface{}) error {
			valuePtr := reflect.ValueOf(v)
			value := valuePtr.Elem()
			value.Set(reflect.ValueOf(nil))

			return nil
		}
	}

	return func(reader io.Reader, v interface{}) error {
		valuePtr := reflect.ValueOf(v)
		value := valuePtr.Elem()

		if !value.CanSet() {
			return errors.New("value can't be set")
		}

		bval, err := readValue(bt, reader)
		if err != nil {
			return err
		}

		converted, err := decodeItem(value.Type(), bt, bval)
		if err != nil {
			return fmt.Errorf("storage can't be converted to type: %w", err)
		}

		value.Set(reflect.ValueOf(converted))

		return nil
	}
}

func convertToType(rt reflect.Type, val interface{}) (interface{}, error) {
	switch rt.Kind() {
	case reflect.Interface:
		return val, nil
	case reflect.Ptr:
		return convertToType(rt.Elem(), val)
	default:
		return convertToKind(rt.Kind(), val)
	}
}

func convertToKind(rk reflect.Kind, v interface{}) (interface{}, error) {
	switch rk {
	case reflect.Int:
		return conv.Int(v)
	case reflect.Int8:
		return conv.Int8(v)
	case reflect.Int16:
		return conv.Int16(v)
	case reflect.Int32:
		return conv.Int32(v)
	case reflect.Int64:
		return conv.Int64(v)
	case reflect.Uint:
		return conv.Uint(v)
	case reflect.Uint8:
		return conv.Uint8(v)
	case reflect.Uint16:
		return conv.Uint16(v)
	case reflect.Uint32:
		return conv.Int32(v)
	case reflect.Uint64:
		return conv.Uint64(v)
	case reflect.Bool:
		return conv.Bool(v)
	case reflect.String:
		return conv.String(v)
	}

	return v, nil
}
