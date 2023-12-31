package bitfield

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPiece(t *testing.T) {
	bf := Bitfield{0b01010100, 0b01010100}
	outputs := []bool{false, true, false, true, false, true, false, false, false, true, false, true, false, true, false, false, false, false, false, false}
	for i := 0; i < len(outputs); i++ {
		assert.Equal(t, outputs[i], bf.HasPiece(i))
	}
}

func TestSetPiece(t *testing.T) {
	tests := []struct {
		input  Bitfield
		output Bitfield
		index  int
	}{
		{
			input:  Bitfield{0b01010100, 0b01010100},
			output: Bitfield{0b01011100, 0b01010100},
			index:  4,
		},
		{
			input:  Bitfield{0b01010100, 0b01010100},
			output: Bitfield{0b01010100, 0b01010100},
			index:  9,
		},
		{
			input:  Bitfield{0b01010100, 0b01010100},
			output: Bitfield{0b01010100, 0b01010100},
			index:  15,
		},
		{
			input:  Bitfield{0b01010100, 0b01010100},
			output: Bitfield{0b01010100, 0b01010100},
			index:  15,
		},
	}
	for _, test := range tests {
		bf := test.input
		bf.SetPiece(test.index)
		assert.Equal(t, test.output, bf)
	}
}
