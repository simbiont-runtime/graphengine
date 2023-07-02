// ---

package resolver

import (
	"testing"

	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewMultiKeysNotifier(t *testing.T) {
	assert := assert.New(t)

	n := NewMultiKeysNotifier(5)

	for i := 0; i < 5; i++ {
		var err error
		if i%2 == 0 {
			err = errors.New("mock error")
		}
		n.Notify(err)
	}

	errs := n.Wait()
	assert.NotNil(errs)
	assert.Equal(3, len(errs))
	assert.Error(errs[2], "mock error")
}
