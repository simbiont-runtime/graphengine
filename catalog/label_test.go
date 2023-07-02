// ---

package catalog

import (
	"testing"

	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/stretchr/testify/assert"
)

func TestNewLabel(t *testing.T) {
	meta := &model.LabelInfo{
		ID:   1,
		Name: model.NewCIStr("test-label"),
	}

	label := NewLabel(meta)
	assert := assert.New(t)
	assert.Equal(label.Meta(), meta)
}
