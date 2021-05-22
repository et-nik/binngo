package encode

import (
	"testing"

	"github.com/et-nik/binngo/binn"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tests := []struct {
		string   string
		expected []byte
	}{
		{
			"test",
			[]byte{binn.StringType, 4, 't', 'e', 's', 't', 0x00},
		},
		{
			"longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglong",
			[]byte{binn.StringType, 0x80, 0x00, 0x00, 0xd0, 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 'l', 'o', 'n', 'g', 0x00},
		},
	}

	for _, test := range tests {
		t.Run(test.string, func(t *testing.T) {
			b, err := Marshal(test.string)

			assert.Nil(t, err)
			assert.Equal(t, test.expected, b)
		})
	}
}
