// ---

package resolver

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/cockroachdb/pebble"
	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/sourcegraph/conc"
	"github.com/twmb/murmur3"
)

// Scheduler is used to schedule Resolve tasks.
type Scheduler struct {
	running   atomic.Bool
	mu        sync.Mutex
	db        *pebble.DB
	size      int
	resolvers []*resolver
	wg        conc.WaitGroup
	cancelFn  context.CancelFunc
}

func NewScheduler(size int) *Scheduler {
	s := &Scheduler{
		size: size,
	}
	return s
}

func (s *Scheduler) SetDB(db *pebble.DB) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db = db
}

// Run initializes the resolvers and start to accept resolve tasks.
func (s *Scheduler) Run() {
	if s.running.Swap(true) {
		return
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	for i := 0; i < s.size; i++ {
		r := newResolver(s.db)
		s.resolvers = append(s.resolvers, r)
		s.wg.Go(func() { r.run(ctx) })
	}
	s.cancelFn = cancelFn
}

// Resolve submits a bundle of keys to resolve
func (s *Scheduler) Resolve(keys []kv.Key, startVer, commitVer kv.Version, notifier Notifier) {
	if len(keys) == 0 {
		return
	}
	for _, key := range keys {
		idx := int(murmur3.Sum32(key)) % s.size
		s.resolvers[idx].push(Task{
			Key:       key,
			StartVer:  startVer,
			CommitVer: commitVer,
			Notifier:  notifier,
		})
	}
}

func (s *Scheduler) Close() {
	s.cancelFn()
	s.wg.Wait()
	for _, r := range s.resolvers {
		r.close()
	}
}
