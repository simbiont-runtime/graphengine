// ---

package catalog

import "github.com/simbiont-runtime/graphengine/parser/model"

// Index represents a runtime index object.
type Index struct {
	meta *model.IndexInfo
}

// NewIndex returns a new index object.
func NewIndex(meta *model.IndexInfo) *Index {
	return &Index{
		meta: meta,
	}
}

// Meta returns the meta information object of this index.
func (i *Index) Meta() *model.IndexInfo {
	return i.meta
}
