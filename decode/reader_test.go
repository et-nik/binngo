package decode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadSize(t *testing.T) {
	tests := []struct {
		name string
		binary []byte
		expectedSize int
		expectedLen int
	}{
		{
			"null",
			[]byte{0x00},
			0,
			1,
		},
		{
			"1 byte size",
			[]byte{0x05},
			5,
			1,
		},
		{
			"4 byte little size",
			[]byte{0x80, 0x00, 0x00, 0x05},
			5,
			4,
		},
		{
			"4 byte size",
			[]byte{0x80, 0x00, 0x01, 0x01},
			257,
			4,
		},
		{
			"4 byte max size",
			[]byte{0xFF, 0xFF, 0xFF, 0xFF},
			2147483647,
			4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sz, ln, err := readSize(bytes.NewReader(test.binary))

			require.NoError(t, err)
			assert.Equal(t, readLen(test.expectedLen), ln)
			assert.Equal(t, test.expectedSize, sz)
		})
	}
}
