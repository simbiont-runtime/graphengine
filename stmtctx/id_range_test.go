// ---

package stmtctx

import (
	"sync/atomic"
	"testing"

	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/assert"
)

func TestIDRange_Next(t *testing.T) {
	rang := NewIDRange(0, 3)

	assert := assert.New(t)

	for i := 0; i < 3; i++ {
		id, err := rang.Next()
		assert.Nil(err)
		assert.Equal(int64(i+1), id)
	}

	_, err := rang.Next()
	assert.NotNil(err)
}

func TestIDRange_NextParallel(t *testing.T) {
	const n = 10000000
	rang := NewIDRange(0, n)

	const p = 20

	wg := conc.WaitGroup{}
	total := atomic.Int64{}
	for i := 0; i < p; i++ {
		wg.Go(func() {
			var subn int64
			for {
				_, err := rang.Next()
				if err != nil {
					total.Add(subn)
					return
				}
				subn++
			}
		})
	}

	wg.Wait()
	assert.Equal(t, total.Load(), int64(n))
}
