// ---

package gc

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/simbiont-runtime/graphengine/storage/resolver"
	"github.com/sourcegraph/conc"
)

// Manager represents the GC manager which is used to scheduler GC tasks to GC worker.
type Manager struct {
	running  atomic.Bool
	size     int
	mu       sync.RWMutex
	db       *pebble.DB
	resolver *resolver.Scheduler
	workers  []*worker
	wg       conc.WaitGroup
	cancelFn context.CancelFunc
	pending  chan Task
}

func NewManager(size int) *Manager {
	return &Manager{
		size:    size,
		pending: make(chan Task, 32),
	}
}

func (m *Manager) SetDB(db *pebble.DB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.db = db
}

func (m *Manager) SetResolver(resolver *resolver.Scheduler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resolver = resolver
}

func (m *Manager) Run() {
	if m.running.Swap(true) {
		return
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	for i := 0; i < m.size; i++ {
		worker := newWorker(m.db, m.resolver)
		m.workers = append(m.workers, worker)
		m.wg.Go(func() { worker.run(ctx, m.pending) })
	}
	m.cancelFn = cancelFn

	// Schedule tasks
	m.wg.Go(func() { m.scheduler(ctx) })
}

func (m *Manager) scheduler(ctx context.Context) {
	const interval = time.Second * 5
	timer := time.NewTimer(interval)
	for {
		select {
		case <-timer.C:
			if len(m.pending) < cap(m.pending) {
				// TODO: schedule some new tasks.
			}

		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) Close() {
	m.cancelFn()
	m.wg.Wait()
}
