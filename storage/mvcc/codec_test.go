// ---

package mvcc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	t1 := Encode([]byte("test"), 1)
	t2 := Encode([]byte("test"), 2)
	r := bytes.Compare(t1, t2)
	assert.True(t, r > 0)
}

func TestDecode(t *testing.T) {
	r, v, err := Decode([]byte("test"))
	assert.ErrorContains(t, err, "insufficient bytes to decode value")
	assert.True(t, v == 0)
	assert.Nil(t, r)

	t1 := Encode([]byte("test"), 1)
	r, v, err = Decode(t1)
	assert.Nil(t, err)
	assert.True(t, v == 1)
	assert.True(t, bytes.Equal([]byte("test"), r))
}
