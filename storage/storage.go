// ---

package storage

import (
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/simbiont-runtime/graphengine/storage/gc"
	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/simbiont-runtime/graphengine/storage/latch"
	"github.com/simbiont-runtime/graphengine/storage/resolver"
)

type mvccStorage struct {
	db        *pebble.DB
	latches   *latch.LatchesScheduler
	resolver  *resolver.Scheduler
	gcManager *gc.Manager
}

// Open returns a new storage instance.
func Open(dirname string, options ...Option) (kv.Storage, error) {
	opt := &pebble.Options{}
	for _, op := range options {
		op(opt)
	}
	db, err := pebble.Open(dirname, opt)
	if err != nil {
		return nil, err
	}

	s := &mvccStorage{
		db:        db,
		latches:   latch.NewScheduler(8),
		resolver:  resolver.NewScheduler(4),
		gcManager: gc.NewManager(2),
	}
	s.resolver.SetDB(db)
	s.gcManager.SetDB(db)
	s.gcManager.SetResolver(s.resolver)

	// Run all background services.
	s.latches.Run()
	s.resolver.Run()
	s.gcManager.Run()

	return s, nil
}

// Begin implements the Storage interface
func (s *mvccStorage) Begin() (kv.Transaction, error) {
	curVer := s.CurrentVersion()
	snap, err := s.Snapshot(curVer)
	if err != nil {
		return nil, err
	}
	txn := &Txn{
		vp:        s,
		db:        s.db,
		us:        NewUnionStore(snap),
		latches:   s.latches,
		resolver:  s.resolver,
		valid:     true,
		startTime: time.Now(),
		startVer:  curVer,
		snapshot:  snap,
	}
	return txn, nil
}

// Snapshot implements the Storage interface.
func (s *mvccStorage) Snapshot(ver kv.Version) (kv.Snapshot, error) {
	snap := &KVSnapshot{
		db:       s.db,
		vp:       s,
		ver:      ver,
		resolver: s.resolver,
	}
	return snap, nil
}

// CurrentVersion implements the VersionProvider interface.
// Currently, we use the system time as our startVer, and the system time
// rewind cannot be tolerant.
func (s *mvccStorage) CurrentVersion() kv.Version {
	return kv.Version(time.Now().UnixNano())
}

// Close implements the Storage interface.
func (s *mvccStorage) Close() error {
	s.latches.Close()
	s.resolver.Close()
	s.gcManager.Close()
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
