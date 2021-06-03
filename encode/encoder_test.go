package encode_test

import (
	"testing"

	"github.com/et-nik/binngo/binn"
	"github.com/et-nik/binngo/encode"
	"github.com/stretchr/testify/assert"
)

func TestEncodeStruct(t *testing.T) {
	v := struct{
		Val1 int64
		Val2 string
	} {
		123,
		"value",
	}

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		binn.ObjectType,
		23,									// total size
		2,									// key/value pairs

		4, 'V', 'a', 'l', '1',				// key
		binn.Uint8Type,						// value type
		123,								// value

		4, 'V', 'a', 'l', '2', 				// key
		binn.StringType,					// value type
		5, 'v', 'a', 'l', 'u', 'e', 0x00, 	// value

	}, result)
}

func TestEncodeMapObjectWithIntValue(t *testing.T) {
	v := map[string]int{"key": 123}

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		binn.ObjectType,
		9,									// total size
		1,									// key/value pairs

		3, 'k', 'e', 'y',					// key
		binn.Uint8Type,						// value type
		123,								// value

	}, result)
}

func TestEncodeMapObjectWithStringValue(t *testing.T) {
	v := map[string]string{"hello": "world"}

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		binn.ObjectType,				// [type] object (container)
		0x11,							// [size] container total size
		0x01,							// [count] key/value pairs
		0x05, 'h', 'e', 'l', 'l', 'o',	// key
		0xA0,							// [type] = string
		0x05,							// [size]
		'w', 'o', 'r', 'l', 'd', 0x00,	// [data] (null terminated)

	}, result)
}

func TestEncodeMap(t *testing.T) {
	v := map[int]int{9: 9}

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		binn.MapType,
		9,									// total size
		1,									// key/value pairs

		0x00, 0x00, 0x00, 0x09,				// key
		binn.Uint8Type,						// value type
		9,									// value

	}, result)
}

func TestEncodeUint(t *testing.T) {
	v := 123

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{binn.Uint8Type, 123}, result)
}

func TestEncodeInt16(t *testing.T) {
	v := -800

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{binn.Int16Type, 0xFC, 0xE0}, result)
}

func TestEncodeUint64(t *testing.T) {
	v := 1099511892096

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{binn.Uint64Type, 0x00, 0x00, 0x01, 0x00, 0x00, 0x04, 0x08, 0x80}, result)
}

func TestEncodeString(t *testing.T) {
	v := "test"

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{binn.StringType, 4, 't', 'e', 's', 't', 0x00}, result)
}

func TestEncodeList(t *testing.T) {
	v := []int{123, -456, 789}

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		binn.ListType,	// [type] list (container)
		0x0B,			// [size] container total size
		0x03,			// [count] items
		0x20, 			// [type] = uint8
		0x7B,			// [data] (123)
		0x41,			// [type] = int16
		0xFE, 0x38,		// [data] (-456)
		0x40,			// [type] = uint16
		0x03, 0x15,		// [data] (789)

	}, result)
}

func TestListInsideMap(t *testing.T) {
	v := map[int][]int{
		2: {-12345, 6789},
	}

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		binn.MapType,				// [type] list (container)
		0x10,						// [size] container total size
		0x01,						// [count] items

		0x00, 0x00, 0x00, 0x02, 	// key

			binn.ListType,			// [type] list (container)
			0x09,					// [size] container total size
			0x02,					// [count] items
			binn.Int16Type,			// [type] = int16
			0xCF, 0xC7,				// [data] (-12345)
			binn.Uint16Type,		// [type] = uint16
			0x1A, 0x85,				// [data] (6789)
	}, result)
}

func TestListInterface(t *testing.T) {
	var v [2]interface{}
	v[0] = 123
	v[1] = "string"

	result, err := encode.Marshal(v)

	assert.Nil(t, err)
	assert.Equal(t, []byte{
		binn.ListType,							// [type] list (container)
		14,										// [size] container total size
		0x02,									// [count] items
		binn.Uint8Type,							// [type] = uint8
		0x7B,									// [data] (123)
		binn.StringType,						// [type] = string
		0x06,									// [size] string len,
		's', 't', 'r', 'i', 'n', 'g', 0x00, 	// [data] null terminated
	}, result)
}

func TestStruct(t *testing.T) {
	// Arrange
	type obj struct {
		Var1 int64          `binn:"object-0"`
		Var2 string         `binn:"object-1"`
		Var3 map[int]string `binn:"object-2-innerMap"`
	}
	val := obj{1099511892096, "string", map[int]string{-20: "innerMap -20"}}

	// Act
	result, err := encode.Marshal(val)

	// Assert
	if assert.Nil(t, err) {
		assert.Equal(t, []byte{
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

		}, result)
	}
}

type custom struct {
	A int
	B string
}

func (c custom) MarshalBINN() ([]byte, error) {
	v := []interface{}{c.A, c.B}
	return encode.Marshal(v)
}

func TestEncodeCustom(t *testing.T) {
	v := custom{500, "custom"}

	result, err := encode.Marshal(v)

	if assert.Nil(t, err) {
		assert.Equal(t, []byte{
			binn.ListType,
			0x0f,									// [size] container total size
			2,										// [count] items
			binn.Uint16Type,						// [type] = uint16
			0x01, 0xf4,								// [data] (500)
			binn.StringType,						// [type] = string
			0x06,									// [size] string len,
			'c', 'u', 's', 't', 'o', 'm', 0x00, 	// [data] null terminated
		}, result)
	}
}
