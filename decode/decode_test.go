package decode

import (
	"bytes"
	"testing"

	"github.com/et-nik/binngo/binn"
	"github.com/stretchr/testify/assert"
)

func TestDecodeIntList(t *testing.T) {
	b := []byte{
		binn.ListType,	// [type] list (container)
		0x0B,			// [size] container total size
		0x03,			// [count] items
		0x20, 			// [type] = uint8
		0x7B,			// [data] (123)
		0x41,			// [type] = int16
		0xFE, 0x38,		// [data] (-456)
		0x40,			// [type] = uint16
		0x03, 0x15,		// [data] (789)
	}
	v := []int{}
	r := bytes.NewReader(b)

	err := decode(r, &v)

	assert.Nil(t, err)
	assert.Equal(t, []int{123, -456, 789}, v)
}

func TestDecodeStringList(t *testing.T) {
	b := []byte{
		binn.ListType,							// [type] list (container)
		23,										// [size] container total size
		0x02,									// [count] items
		binn.StringType, 						// [type] = string
		0x05,									// [size]
		'h', 'e', 'l', 'l', 'o', 0x00,			// [data] (null terminated)
		binn.StringType, 						// [type] = string
		0x05,									// [size]
		'w', 'o', 'r', 'l', 'd', 0x00,			// [data] (null terminated)
	}
	v := []string{}
	r := bytes.NewReader(b)

	err := decode(r, &v)

	assert.Nil(t, err)
	assert.Equal(t, []string{"hello", "world"}, v)
}

func TestDecodeIterfaceList(t *testing.T) {
	b := []byte{
		binn.ListType,							// [type] list (container)
		26,										// [size] container total size
		0x03,									// [count] items
		binn.StringType, 						// [type] = string
		0x05,									// [size]
		'h', 'e', 'l', 'l', 'o', 0x00,			// [data] (null terminated)
		binn.StringType, 						// [type] = string
		0x05,									// [size]
		'w', 'o', 'r', 'l', 'd', 0x00,			// [data] (null terminated)
		0x40,								    // [type] = uint16
		0x03, 0x15,								// [data] (789)
	}
	var v []interface{}
	r := bytes.NewReader(b)

	err := decode(r, &v)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"hello", "world", uint16(789)}, v)
}

func TestDecodeStringMap(t *testing.T) {
	b := []byte{
		0xE1,								// [type] map (container)
		0x1A,								// [size] container total size
		0x02,								// [count] key/value pairs
		0x00, 0x00, 0x00, 0x01, 			// key
		0xA0,								// [type] = string
		0x03,             					// [size]
		'a', 'd', 'd', 0x00,          		// [data] (null terminated)
		0x00, 0x00, 0x00, 0x02,				// key
		0xE0,             					// [type] list (container)
		0x09,             					// [size] container total size
		0x02,             					// [count] items
		0x41,             					// [type] = int16
		0xCF, 0xC7,         				// [data] (-12345)
		0x40,             					// [type] = uint16
		0x1A, 0x85,         				// [data] (6789)
	}
	v := map[int]interface{}{}
	r := bytes.NewReader(b)

	err := decode(r, &v)

	assert.Nil(t, err)
	assert.Equal(t, map[int]interface{}{
		1: "add",
		2: []interface{}{int16(-12345), uint16(6789)},
	}, v)
}

func TestDecodeStringObjectToStruct(t *testing.T) {
	b := []byte{
		0xE2,									// type = object (container)
		0x11,									// container total size
		0x01,									// key/value pairs count
		0x05, 'h', 'e', 'l', 'l', 'o',			// key
		0xA0,									// type = string
		0x05, 'w', 'o', 'r', 'l', 'd', 0x00,	// value (null terminated)
	}
	type ts struct {
		Hello string `binn:"hello"`
	}
	var v ts
	r := bytes.NewReader(b)

	err := decode(r, &v)

	assert.Nil(t, err)
	assert.Equal(t, ts{Hello: "world"}, v)
}

func TestDecodeStringObjectToMap(t *testing.T) {
	b := []byte{
		0xE2,									// type = object (container)
		0x11,									// container total size
		0x01,									// key/value pairs count
		0x05, 'h', 'e', 'l', 'l', 'o',			// key
		0xA0,									// type = string
		0x05, 'w', 'o', 'r', 'l', 'd', 0x00,	// value (null terminated)
	}
	m := map[string]string{}
	r := bytes.NewReader(b)

	err := decode(r, &m)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"hello": "world"}, m)
}
