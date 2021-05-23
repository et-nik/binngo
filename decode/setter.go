package decode

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/et-nik/binngo/binn"
)

var kindMapper = map[binn.Type]reflect.Kind{
	binn.Int8Type:   reflect.Int8,
	binn.Int16Type:  reflect.Int16,
	binn.Int32Type:  reflect.Int32,
	binn.Int64Type:  reflect.Int64,
	binn.Uint8Type:  reflect.Uint8,
	binn.Uint16Type: reflect.Uint16,
	binn.Uint32Type: reflect.Uint32,
	binn.Uint64Type: reflect.Uint64,
	binn.StringType: reflect.String,
}

func addSliceItem(btype binn.Type, bval []byte, v interface{}) error {
	value := reflect.ValueOf(v).Elem()

	var err error

	if value.Kind() != reflect.Slice {
		return &UnknownValueError{reflect.Slice, value.Kind()}
	}

	val, err := decodeItem(value.Type().Elem(), btype, bval)

	if err != nil {
		return err
	}

	if !value.CanSet() {
		return ErrCantSetValue
	}

	value.Set(
		reflect.Append(value, reflect.Indirect(reflect.ValueOf(val))),
	)

	return nil
}

func addMapItem(k interface{}, bt binn.Type, bval []byte, v interface{}) error {
	value := reflect.ValueOf(v).Elem()

	var err error

	kind := value.Kind()

	if kind != reflect.Map {
		return &UnknownValueError{reflect.Map, value.Kind()}
	}

	val, err := decodeItem(value.Type().Elem(), bt, bval)

	if err != nil {
		return fmt.Errorf("failed to add map item: %w", err)
	}

	if !value.CanSet() {
		return ErrCantSetValue
	}

	value.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(val))

	return nil
}

func addObjectItem(key string, btype binn.Type, bval []byte, v interface{}) error {
	kind := reflect.ValueOf(v).Elem().Kind()

	if kind == reflect.Interface {
		kind = reflect.ValueOf(v).Elem().Elem().Kind()
	}

	switch kind {
	case reflect.Map:
		return addMapItem(key, btype, bval, v)
	case reflect.Struct:
		return addObjectItemToStruct(key, btype, bval, v)
	case reflect.Ptr:
		return addObjectItem(key, btype, bval, reflect.ValueOf(v).Elem().Interface())
	}

	return nil
}

func addObjectItemToStruct(k string, bt binn.Type, bval []byte, v interface{}) error {
	value := reflect.ValueOf(v).Elem()

	if value.Kind() == reflect.Interface {
		value = reflect.Indirect(value.Elem())
	}

	if value.Kind() != reflect.Struct {
		return &UnknownValueError{reflect.Struct, value.Kind()}
	}

	field := value.FieldByName(k)

	if !field.IsValid() {
		fieldName, err := findFieldNameByTag(k, value.Type())
		if err != nil {
			return fmt.Errorf("failed to find field name by tag: %w", err)
		}

		field = reflect.Indirect(value.FieldByName(fieldName))

		if !field.IsValid() {
			return ErrInvalidStructValue
		}
	}

	var err error

	val, err := decodeItem(field.Type(), bt, bval)

	if err != nil {
		return fmt.Errorf("failed to add object item to struct: %w", err)
	}

	if !field.CanSet() {
		return ErrCantSetValue
	}

	field.Set(reflect.ValueOf(val))

	return nil
}

func findFieldNameByTag(key string, rt reflect.Type) (string, error) {
	if rt.Kind() != reflect.Struct {
		return "", ErrInvalidItem
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get("binn"), ",")[0]
		if v == key {
			return f.Name, nil
		}
	}

	return "", ErrItemNotFound
}
