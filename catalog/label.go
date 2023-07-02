// ---

package catalog

import (
	"sync"

	"github.com/simbiont-runtime/graphengine/parser/model"
)

// Label represents a runtime label object.
type Label struct {
	mu sync.RWMutex

	meta *model.LabelInfo
}

// NewLabel returns a label instance.
func NewLabel(meta *model.LabelInfo) *Label {
	l := &Label{
		meta: meta,
	}
	return l
}

// Meta returns the meta information object of this label.
func (l *Label) Meta() *model.LabelInfo {
	return l.meta
}
