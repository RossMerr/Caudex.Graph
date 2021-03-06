package expressions

import (
	"testing"
)

func Test_XORSwap(t *testing.T) {

	var tests = []struct {
		a      bool
		result bool
	}{
		{
			a:      false,
			result: true,
		},
		{
			a:      true,
			result: false,
		},
	}

	for i, tt := range tests {
		result := XORSwap(tt.a)
		if result != tt.result {
			t.Errorf("%d. xor:\n  exp=%v\n  got=%v\n\n", i, tt.result, result)
		}
	}
}

func Test_XOR(t *testing.T) {

	var tests = []struct {
		a      bool
		b      bool
		result bool
	}{
		{
			a:      false,
			b:      false,
			result: false,
		},
		{
			a:      false,
			b:      true,
			result: true,
		},
		{
			a:      true,
			b:      false,
			result: true,
		},
		{
			a:      true,
			b:      true,
			result: false,
		},
	}

	for i, tt := range tests {
		result := XORExclusive(tt.a, tt.b)
		if result != tt.result {
			t.Errorf("%d. xor:\n  exp=%v\n  got=%v\n\n", i, tt.result, result)
		}
	}
}
