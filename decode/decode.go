package decode

import (
	"bytes"
	"io"
	"reflect"
	"sync"

	"github.com/cstockton/go-conv"
	"github.com/et-nik/binngo/binn"
	"github.com/et-nik/binngo/encode"
)

const (
	maxOneByteSize = 127
)

type decodeFunc func(reader io.Reader, v interface{}) error

var decoderCache sync.Map // map[binn.Type]decodeFunc

type readLen int

func decode(reader io.Reader, v interface{}) error {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{rt}
	}

	containerType, _, err := readType(reader)
	if err != nil {
		return err
	}

	if rt.Implements(unmarshalerType) && isStorageContainer(containerType) {
		return decodeUnmarshalerStorage(containerType, reader, v)
	}

	return decodeStorage(containerType, reader, v)
}

func decodeUnmarshalerStorage(containerType binn.Type, reader io.Reader, v interface{}) error {
	size, ln, err := readSize(reader)
	if err != nil {
		return err
	}

	typeSize := len(encode.Int(int(containerType)))

	buf := make([]byte, size-int(ln)-typeSize)
	n, err := reader.Read(buf)
	if err != nil {
		return err
	}

	if n != len(buf) {
		return ErrIncompleteRead
	}

	data := make([]byte, 0, size)
	data = append(data, encode.Uint8(uint8(containerType))...)
	data = append(data, encode.Size(size, false)...)
	data = append(data, buf...)

	err = decodeUnmarshaler(data, v)
	if err != nil {
		return err
	}

	return nil
}

func decodeStorage(containerType binn.Type, reader io.Reader, v interface{}) error {
	decoder := loadDecodeFunc(containerType)
	return decoder(reader, v)
}

//nolint:funlen
func decodeItem(rt reflect.Type, btype binn.Type, bval []byte) (interface{}, error) {
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
		v = Uint8(bval)
	case binn.Uint16Type:
		v = Uint16(bval)
	case binn.Uint32Type:
		v = Uint32(bval)
	case binn.Uint64Type:
		v = Uint64(bval)
	case binn.Int8Type:
		v = Int8(bval)
	case binn.Int16Type:
		v = Int16(bval)
	case binn.Int32Type:
		v = Int32(bval)
	case binn.Int64Type:
		v = Int64(bval)
	case binn.Float32Type:
		v = Float32(bval)
	case binn.Float64Type:
		v = Float64(bval)
	case binn.StringType:
		v = String(bval[:len(bval)-1])
	case binn.BlobType:
		v = bval
	case binn.ListType:
		var l []interface{}
		br := bytes.NewReader(bval)
		_, wasReadLen, _ := readSize(br)
		cnt, wasReadCnt, _ := readSize(br)
		wasRead := wasReadLen + wasReadCnt

		err = decodeListItems(br, &l, len(bval)-int(wasRead), wasRead, cnt)
		if err != nil {
			return nil, err
		}

		return l, nil
	case binn.MapType:
		br := bytes.NewReader(bval)
		sz, rlsize, _ := readSize(br)
		cnt, rlcnt, _ := readSize(br)

		mapType := reflect.MapOf(reflect.TypeOf(int(0)), rt.Elem())
		ptr := reflect.New(mapType)
		ptr.Elem().Set(reflect.MakeMap(mapType))
		m := ptr.Interface()

		err = decodeMapItems(br, m, sz, rlsize+rlcnt, cnt)
		if err != nil {
			return nil, err
		}

		return reflect.ValueOf(m).Elem().Interface(), nil
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

func loadDecodeFunc(bt binn.Type) decodeFunc {
	if fi, ok := decoderCache.Load(bt); ok {
		return fi.(decodeFunc)
	}

	var (
		wg sync.WaitGroup
		f  decodeFunc
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

func newTypeDecoder(bt binn.Type) decodeFunc {
	switch bt {
	case binn.ListType:
		return decodeList
	case binn.MapType:
		return decodeMap
	case binn.ObjectType:
		return decodeObject
	case binn.Null:
		return func(_ io.Reader, _ interface{}) error {
			return nil
		}
	}

	decoder := newValueDecoder(bt)
	return decoder.DecodeValue
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
		return conv.Uint32(v)
	case reflect.Uint64:
		return conv.Uint64(v)
	case reflect.Bool:
		return conv.Bool(v)
	case reflect.String:
		return conv.String(v)
	}

	return v, nil
}
