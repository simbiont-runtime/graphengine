// ---

package codec

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFastSlowFastReverse(t *testing.T) {
	if !supportsUnaligned {
		return
	}
	b := []byte{1, 2, 3, 4, 5, 6, 7, 8, 255, 0, 0, 0, 0, 0, 0, 0, 0, 247}
	r1 := b
	fastReverseBytes(b)
	r2 := b
	reverseBytes(r2)
	require.Equal(t, r1, r2)
}

func TestBytesCodec(t *testing.T) {
	inputs := []struct {
		enc  []byte
		dec  []byte
		desc bool
	}{
		{[]byte{}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 247}, false},
		{[]byte{}, []byte{255, 255, 255, 255, 255, 255, 255, 255, 8}, true},
		{[]byte{0}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 248}, false},
		{[]byte{0}, []byte{255, 255, 255, 255, 255, 255, 255, 255, 7}, true},
		{[]byte{1, 2, 3}, []byte{1, 2, 3, 0, 0, 0, 0, 0, 250}, false},
		{[]byte{1, 2, 3}, []byte{254, 253, 252, 255, 255, 255, 255, 255, 5}, true},
		{[]byte{1, 2, 3, 0}, []byte{1, 2, 3, 0, 0, 0, 0, 0, 251}, false},
		{[]byte{1, 2, 3, 0}, []byte{254, 253, 252, 255, 255, 255, 255, 255, 4}, true},
		{[]byte{1, 2, 3, 4, 5, 6, 7}, []byte{1, 2, 3, 4, 5, 6, 7, 0, 254}, false},
		{[]byte{1, 2, 3, 4, 5, 6, 7}, []byte{254, 253, 252, 251, 250, 249, 248, 255, 1}, true},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 247}, false},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0}, []byte{255, 255, 255, 255, 255, 255, 255, 255, 0, 255, 255, 255, 255, 255, 255, 255, 255, 8}, true},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 255, 0, 0, 0, 0, 0, 0, 0, 0, 247}, false},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8}, []byte{254, 253, 252, 251, 250, 249, 248, 247, 0, 255, 255, 255, 255, 255, 255, 255, 255, 8}, true},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 255, 9, 0, 0, 0, 0, 0, 0, 0, 248}, false},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, []byte{254, 253, 252, 251, 250, 249, 248, 247, 0, 246, 255, 255, 255, 255, 255, 255, 255, 7}, true},
	}

	for _, input := range inputs {
		require.Len(t, input.dec, EncodedBytesLength(len(input.enc)))

		if input.desc {
			b := EncodeBytesDesc(nil, input.enc)
			require.Equal(t, input.dec, b)

			_, d, err := DecodeBytesDesc(b, nil)
			require.NoError(t, err)
			require.Equal(t, input.enc, d)
		} else {
			b := EncodeBytes(nil, input.enc)
			require.Equal(t, input.dec, b)

			_, d, err := DecodeBytes(b, nil)
			require.NoError(t, err)
			require.Equal(t, input.enc, d)
		}
	}

	// Test error decode.
	errInputs := [][]byte{
		{1, 2, 3, 4},
		{0, 0, 0, 0, 0, 0, 0, 247},
		{0, 0, 0, 0, 0, 0, 0, 0, 246},
		{0, 0, 0, 0, 0, 0, 0, 1, 247},
		{1, 2, 3, 4, 5, 6, 7, 8, 0},
		{1, 2, 3, 4, 5, 6, 7, 8, 255, 1},
		{1, 2, 3, 4, 5, 6, 7, 8, 255, 1, 2, 3, 4, 5, 6, 7, 8},
		{1, 2, 3, 4, 5, 6, 7, 8, 255, 1, 2, 3, 4, 5, 6, 7, 8, 255},
		{1, 2, 3, 4, 5, 6, 7, 8, 255, 1, 2, 3, 4, 5, 6, 7, 8, 0},
	}

	for _, input := range errInputs {
		_, _, err := DecodeBytes(input, nil)
		require.Error(t, err)
	}
}

func TestBytesCodecExt(t *testing.T) {
	inputs := []struct {
		enc []byte
		dec []byte
	}{
		{[]byte{}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 247}},
		{[]byte{1, 2, 3}, []byte{1, 2, 3, 0, 0, 0, 0, 0, 250}},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 255, 9, 0, 0, 0, 0, 0, 0, 0, 248}},
	}

	// `assertEqual` is to deal with test case for `[]byte{}` & `[]byte(nil)`.
	assertEqual := func(expected []byte, acutal []byte) {
		require.Equal(t, len(expected), len(acutal))
		for i := range expected {
			require.Equal(t, expected[i], acutal[i])
		}
	}

	for _, input := range inputs {
		assertEqual(input.enc, EncodeBytesExt(nil, input.enc, true))
		assertEqual(input.dec, EncodeBytesExt(nil, input.enc, false))
	}
}

func TestFloatCodec(t *testing.T) {
	tblFloat := []float64{
		-1,
		0,
		1,
		math.MaxFloat64,
		math.MaxFloat32,
		math.SmallestNonzeroFloat32,
		math.SmallestNonzeroFloat64,
		math.Inf(-1),
		math.Inf(1),
	}

	for _, floatNum := range tblFloat {
		b := EncodeFloat(nil, floatNum)
		_, v, err := DecodeFloat(b)
		require.NoError(t, err)
		require.Equal(t, floatNum, v)

		b = EncodeFloatDesc(nil, floatNum)
		_, v, err = DecodeFloatDesc(b)
		require.NoError(t, err)
		require.Equal(t, floatNum, v)
	}

	tblCmp := []struct {
		Arg1 float64
		Arg2 float64
		Ret  int
	}{
		{1, -1, 1},
		{1, 0, 1},
		{0, -1, 1},
		{0, 0, 0},
		{math.MaxFloat64, 1, 1},
		{math.MaxFloat32, math.MaxFloat64, -1},
		{math.MaxFloat64, 0, 1},
		{math.MaxFloat64, math.SmallestNonzeroFloat64, 1},
		{math.Inf(-1), 0, -1},
		{math.Inf(1), 0, 1},
		{math.Inf(-1), math.Inf(1), -1},
	}

	for _, floatNums := range tblCmp {
		b1 := EncodeFloat(nil, floatNums.Arg1)
		b2 := EncodeFloat(nil, floatNums.Arg2)

		ret := bytes.Compare(b1, b2)
		require.Equal(t, floatNums.Ret, ret)

		b1 = EncodeFloatDesc(nil, floatNums.Arg1)
		b2 = EncodeFloatDesc(nil, floatNums.Arg2)

		ret = bytes.Compare(b1, b2)
		require.Equal(t, -floatNums.Ret, ret)
	}
}
