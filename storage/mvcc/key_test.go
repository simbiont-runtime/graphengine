// ---

package mvcc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKey(t *testing.T) {
	x := []byte("test")
	ex := NewKey(x)
	k, v, err := Decode(ex)
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(x, k))
	assert.Zero(t, v)
}
