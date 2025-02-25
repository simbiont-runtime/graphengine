// ---

package stmtctx

import (
	"strings"
	"sync"
	"sync/atomic"

	"github.com/simbiont-runtime/graphengine/catalog"
	"github.com/simbiont-runtime/graphengine/storage/kv"
)

// Context represent the intermediate state of a query execution and will be
// reset after a query finished.
type Context struct {
	store   kv.Storage
	catalog *catalog.Catalog

	mu struct {
		sync.RWMutex

		currentGraph string
		txn          *LazyTxn

		affectedRows uint64
		foundRows    uint64
		records      uint64
		deleted      uint64
		updated      uint64
		copied       uint64
		touched      uint64

		warnings   []SQLWarn
		errorCount uint16
	}

	// TODO: perhaps we can move these to a separate struct.
	planID       atomic.Int64
	planColumnID atomic.Int64
}

// New returns a session statement context instance.
func New(store kv.Storage, catalog *catalog.Catalog) *Context {
	return &Context{
		store:   store,
		catalog: catalog,
	}
}

// Reset resets all variables associated to execute a query.
func (sc *Context) Reset() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.mu.affectedRows = 0
	sc.mu.foundRows = 0
	sc.mu.records = 0
	sc.mu.deleted = 0
	sc.mu.updated = 0
	sc.mu.copied = 0
	sc.mu.touched = 0
	sc.mu.warnings = sc.mu.warnings[:0]
	sc.mu.errorCount = 0
}

// Store returns the storage instance.
func (sc *Context) Store() kv.Storage {
	return sc.store
}

// Catalog returns the catalog object.
func (sc *Context) Catalog() *catalog.Catalog {
	return sc.catalog
}

// CurrentGraph returns the current chosen catalog graph
func (sc *Context) CurrentGraph() *catalog.Graph {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.Catalog().Graph(sc.mu.currentGraph)
}

// SetCurrentGraphName changes the current graph name.
func (sc *Context) SetCurrentGraphName(graphName string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.mu.currentGraph = strings.ToLower(graphName)
}

// CurrentGraphName returns the current chosen graph name.
func (sc *Context) CurrentGraphName() string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.mu.currentGraph
}

func (sc *Context) AllocPlanID() int {
	return int(sc.planID.Add(1))
}

func (sc *Context) AllocPlanColumnID() int64 {
	return sc.planColumnID.Add(1)
}
