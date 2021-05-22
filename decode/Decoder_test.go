package decode_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/et-nik/binngo/binn"
	"github.com/et-nik/binngo/decode"

	"github.com/stretchr/testify/assert"
)

func TestUnknownValueError_ExpectedBool(t *testing.T) {
	b := []byte{binn.True}
	var r int

	err := decode.Unmarshal(b, &r)

	assert.NotNil(t, err)
	var e *decode.UnknownValueError
	assert.ErrorAs(t, err, &e)
	assert.Equal(t, reflect.Bool, err.(*decode.UnknownValueError).Expected)
	assert.Equal(t, reflect.Int, err.(*decode.UnknownValueError).Got)
}

func TestUnknownValueError_ExpectedSlice(t *testing.T) {
	b := []byte{
		binn.ListType,	// [type] list (container)
		0x05,			// [size] container total size
		0x01,			// [count] items
		0x20, 			// [type] = uint8
		0x7B,			// [data] (123)
	}
	var r int

	err := decode.Unmarshal(b, &r)

	assert.NotNil(t, err)
	var e *decode.UnknownValueError
	assert.ErrorAs(t, err, &e)
	assert.Equal(t, reflect.Slice, err.(*decode.UnknownValueError).Expected)
	assert.Equal(t, reflect.Int, err.(*decode.UnknownValueError).Got)
}

func TestInvalidCount(t *testing.T) {
	b := []byte{
		binn.ListType,	// [type] list (container)
		0x05,			// [size] container total size
		0x03,			// [count] items
		0x20, 			// [type] = uint8
		0x7B,			// [data] (123)
	}
	v := []int{}

	err := decode.Unmarshal(b, &v)

	assert.NotNil(t, err)
	assert.Equal(t, []int{123}, v)
}

func TestSimpleStorages(t *testing.T) {
	tests := []struct {
		name      string
		binary    []byte
		expected  interface{}
	}{
		{
			"true",
			[]byte{binn.True},
			true,
		},
		{
			"uint8",
			[]byte{binn.Uint8Type, 33},
			uint8(33),
		},
		{
			"int16",
			[]byte{binn.Int16Type, 0xCF, 0xC7},
			int16(-12345),
		},
		{
			"int32",
			[]byte{binn.Int32Type, 0xFF, 0x43, 0x9E, 0xB2},
			int32(-12345678),
		},
		{
			"uint64",
			[]byte{binn.Int64Type, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE},
			int64(9223372036854775806),
		},
		{
			"string",
			[]byte{binn.StringType, 0x05, 'h', 'e', 'l', 'l', 'o', 0x00},
			"hello",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v interface{}
			err := decode.Unmarshal(test.binary, &v)

			assert.Nil(t, err)
			assert.Equal(t, test.expected, v)
		})
	}
}

func TestUnmashalObject(t *testing.T) {
	b := []byte {
		0xE2,           			// [type] object (container)
		0x14,           			// [size] container total size
		0x02,           			// [count] key/value pairs

		0x02, 'i', 'd',     		// key
		0x20,           			// [type] = uint8
		0x01,           			// [data] (1)

		0x04, 'n', 'a', 'm', 'e',   // key
		0xA0,           			// [type] = string
		0x04,           			// [size]
		'J', 'o', 'h', 'n', 0x00,   // [data] (null terminated)
	}
	type ts struct {
		ID   uint8  `binn:"id"`
		Name string `binn:"name"`
	}
	obj := ts{}

	err := decode.Unmarshal(b, &obj)

	assert.Nil(t, err)
	assert.Equal(t, ts{1, "John"}, obj)
}

func TestUnmarshalStringObjectToMap(t *testing.T) {
	b := []byte{
		0xE2,                          // type = object (container)
		0x11,                          // container total size
		0x01,                          // key/value pairs count
		0x05, 'h', 'e', 'l', 'l', 'o', // key
		0xA0,                                // type = string
		0x05, 'w', 'o', 'r', 'l', 'd', 0x00, // value (null terminated)
	}
	m := map[string]string{}

	err := decode.Unmarshal(b, &m)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"hello": "world"}, m)
}

// https://github.com/liteserver/binn/blob/master/spec.md#a-list-of-objects
func TestListOfObjects(t *testing.T) {
	b := []byte{
		0xE0,           			// [type] list (container)
		0x2B,           			// [size] container total size
		0x02,           			// [count] items

		0xE2,           			// [type] object (container)
		0x14,           			// [size] container total size
		0x02,           			// [count] key/value pairs

		0x02, 'i', 'd',     		// key
		0x20,           			// [type] = uint8
		0x01,           			// [data] (1)

		0x04, 'n', 'a', 'm', 'e',   // key
		0xA0,           			// [type] = string
		0x04,           			// [size]
		'J', 'o', 'h', 'n', 0x00,   // [data] (null terminated)

		0xE2,           			// [type] object (container)
		0x14,           			// [size] container total size
		0x02,          	 			// [count] key/value pairs

		0x02, 'i', 'd',         	// key
		0x20,           			// [type] = uint8
		0x02,           			// [data] (2)

		0x04, 'n', 'a', 'm', 'e',   // key
		0xA0,           			// [type] = string
		0x04,           			// [size]
		'E', 'r', 'i', 'c', 0x00,   // [data] (null terminated)
	}
	type ts struct {
		ID   uint8  `binn:"id"`
		Name string `binn:"name"`
	}
	l := []ts{}

	err := decode.Unmarshal(b, &l)

	if assert.Nil(t, err) {
		assert.Equal(t, []ts{
			{1, "John"},
			{2, "Eric"},
		}, l)
	}
}

func TestJSON(t *testing.T) {
	b := []byte{'t', 'r', 'u', 'e'}
	var r bool

	err := json.Unmarshal(b, &r)

	assert.Nil(t, err)
	assert.True(t, r)
}
