// ---

package stmtctx

import (
	"sync/atomic"

	"github.com/simbiont-runtime/graphengine/catalog"
	"github.com/simbiont-runtime/graphengine/meta"
	"github.com/simbiont-runtime/graphengine/storage/kv"
)

// IDRange represents an ID range. The ID range will be (base, max]
type IDRange struct {
	base atomic.Int64
	max  int64
}

// NewIDRange returns a new ID range.
func NewIDRange(base, max int64) *IDRange {
	idr := &IDRange{
		max: max,
	}
	idr.base.Store(base)
	return idr
}

// Next retrieves the next available ID.
func (r *IDRange) Next() (int64, error) {
	next := r.base.Add(1)
	if next > r.max {
		return 0, ErrIDExhaust
	}
	return next, nil
}

// AllocID allocates n IDs.
func (sc *Context) AllocID(graph *catalog.Graph, n int) (*IDRange, error) {
	graph.MDLock()
	defer graph.MDUnlock()

	var idRange *IDRange
	err := kv.Txn(sc.store, func(txn kv.Transaction) error {
		meta := meta.New(txn)
		base, err := meta.AdvanceID(graph.Meta().ID, n)
		if err != nil {
			return err
		}
		idRange = NewIDRange(base, base+int64(n))
		return nil
	})

	return idRange, err
}
