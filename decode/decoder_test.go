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
			"nil",
			[]byte{binn.Null},
			nil,
		},
		{
			"true",
			[]byte{binn.True},
			true,
		},
		{
			"false",
			[]byte{binn.False},
			false,
		},
		{
			"uint8",
			[]byte{binn.Uint8Type, 33},
			uint8(33),
		},
		{
			"int8",
			[]byte{binn.Int8Type, 0xDF},
			int8(-33),
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
			"uint32",
			[]byte{binn.Uint32Type, 0x00, 0xBC, 0x61, 0x4E},
			int32(12345678),
		},
		{
			"uint64",
			[]byte{binn.Int64Type, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE},
			int64(9223372036854775806),
		},
		{
			"float32",
			[]byte{binn.Float32Type, 0x41, 0x82, 0xCA, 0xC1},
			float32(16.349),
		},
		{
			"float64",
			[]byte{binn.Float64Type, 0x40, 0x30, 0x59, 0x96, 0x65, 0xF5, 0x11, 0x6B},
			16.349951145487847,
		},
		{
			"string",
			[]byte{binn.StringType, 0x05, 'h', 'e', 'l', 'l', 'o', 0x00},
			"hello",
		},
		{
			"blob",
			[]byte{binn.BlobType, 0x05, 0x00, 0x01, 0x02, 0x03, 0x04},
			[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
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

func TestMapInObjectStruct(t *testing.T) {
	b := []byte{
		0xe2,																			// [type] object
		0x4f,																			// [size]
		0x03, 																			// [count]

		0x08, 'o', 'b', 'j', 'e', 'c', 't', '-', '0',  									// [key]
		0x80, 0x00, 0x00, 0x01, 0x00, 0x00, 0x04, 0x08, 0x80,							// [value] (1099511892096)

		0x08, 'o', 'b', 'j', 'e', 'c', 't', '-', '1',  									// key
		0xa0, 0x06, 's', 't', 'r', 'i', 'n', 'g', 0x00,									// [type] string, [value]

		0x11, 'o', 'b', 'j', 'e', 'c', 't', '-', '2', '-', 'i', 'n', 'n', 'e', 'r', 'M', 'a', 'p',
		0xE1, 0x16, 0x01, 																// [type] map, [size], [count]
		0xff, 0xff, 0xff, 0xec, 														// [key] -20
		0xa0, 0x0c, 'i', 'n', 'n', 'e', 'r', 'M', 'a', 'p', ' ', '-', '2', '0', 0x00,

	}
	type obj struct {
		Var1 int64          `binn:"object-0"`
		Var2 string         `binn:"object-1"`
		Var3 map[int]string `binn:"object-2-innerMap"`
	}
	var v obj

	err := decode.Unmarshal(b, &v)

	if assert.Nil(t, err) {
		assert.Equal(t, obj{
			1099511892096,
			"string",
			map[int]string{-20: "innerMap -20"},
		}, v)
	}
}
