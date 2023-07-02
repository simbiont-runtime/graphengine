// ---

package catalog

import (
	"testing"

	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/stretchr/testify/assert"
)

func TestNewIndex(t *testing.T) {
	meta := &model.IndexInfo{
		ID:   1,
		Name: model.NewCIStr("test-index"),
	}

	index := NewIndex(meta)
	assert := assert.New(t)
	assert.Equal(index.Meta(), meta)
}
