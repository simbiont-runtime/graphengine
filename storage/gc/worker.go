// ---

package gc

import (
	"context"
	"sync/atomic"

	"github.com/cockroachdb/pebble"
	"github.com/simbiont-runtime/graphengine/storage/resolver"
)

// worker represents a GC worker which is used to clean staled versions.
type worker struct {
	db       *pebble.DB
	resolver *resolver.Scheduler
	closed   atomic.Bool
}

func newWorker(db *pebble.DB, resolver *resolver.Scheduler) *worker {
	return &worker{
		db:       db,
		resolver: resolver,
	}
}

func (w *worker) run(ctx context.Context, queue <-chan Task) {
	for {
		select {
		case task, ok := <-queue:
			if !ok {
				return
			}
			w.execute(task)

		case <-ctx.Done():
			return
		}
	}
}

func (w *worker) execute(task Task) {

}
