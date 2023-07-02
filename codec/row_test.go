// ---

package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowBytes(t *testing.T) {
	rb := &rowBytes{
		labelIDs:    []uint16{1, 2, 3, 4},
		propertyIDs: []uint16{1, 2, 3, 4},
		offsets:     []uint16{1, 2, 3, 4},
		data:        []byte("abcd"),
	}
	bytes := rb.toBytes(nil)
	rb2 := &rowBytes{}
	err := rb2.fromBytes(bytes)
	assert.NoError(t, err)
	assert.Equal(t, rb, rb2)

	a := rb2.getData(1)
	assert.Equal(t, "b", string(a))

	idx := rb2.findProperty(3)
	assert.Equal(t, 2, idx)
	idx = rb2.findProperty(5)
	assert.Equal(t, -1, idx)
}
