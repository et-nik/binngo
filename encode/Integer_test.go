package encode

import (
	"testing"

	"github.com/et-nik/binngo/binn"
	"github.com/stretchr/testify/assert"
)

func TestIntPack(t *testing.T) {
	tests := []struct {
		name         string
		int          int
		sizeExpected int
		binExpected  []byte
	}{
		{
			"int8 compressed",
			123,
			2,
			[]byte{binn.Uint8Type, 123},
		},
		{
			"int16 compressed",
			789,
			3,
			[]byte{binn.Uint16Type, 0x03, 0x15},
		},
		{
			"int16",
			-12345,
			3,
			[]byte{binn.Int16Type, 0xcf, 0xc7},
		},
		{
			"int32",
			-12345678,
			5,
			[]byte{binn.Int32Type, 0xff, 0x43, 0x9e, 0xb2},
		},
		{
			"uint64",
			9223372036854775806,
			9,
			[]byte{binn.Uint64Type, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := Marshal(test.int)

			assert.Nil(t, err)
			assert.Equal(t, test.binExpected, b)
		})
	}
}
