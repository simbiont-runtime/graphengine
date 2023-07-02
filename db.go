package graphengine

import (
	"sync"

	"github.com/simbiont-runtime/graphengine/catalog"
	"github.com/simbiont-runtime/graphengine/session"
	"github.com/simbiont-runtime/graphengine/storage"
	"github.com/simbiont-runtime/graphengine/storage/kv"
)

// DB represents the  GraphEngine database instance.
type DB struct {
	// All fields are not been protected by Mutex will be read-only.
	options *Options
	store   kv.Storage
	catalog *catalog.Catalog

	mu struct {
		sync.RWMutex
		sessions map[int64]*session.Session
	}
}

// Open opens a  GraphEngine database instance with specified directory name.
func Open(dirname string, opt *Options) (*DB, error) {
	if opt == nil {
		opt = &Options{}
	}
	opt.SetDefaults()

	store, err := storage.Open(dirname)
	if err != nil {
		return nil, err
	}

	// Load the catalog from storage.
	snapshot, err := store.Snapshot(store.CurrentVersion())
	if err != nil {
		return nil, err
	}
	catalog, err := catalog.Load(snapshot)
	if err != nil {
		return nil, err
	}

	db := &DB{
		options: opt,
		store:   store,
		catalog: catalog,
	}
	db.mu.sessions = map[int64]*session.Session{}

	return db, nil
}

// Store returns the storage engine object.
func (db *DB) Store() kv.Storage {
	return db.store
}

// Catalog returns the catalog object.
func (db *DB) Catalog() *catalog.Catalog {
	return db.catalog
}

// NewSession returns a new session.
func (db *DB) NewSession() *session.Session {
	// TODO: concurrency limitation
	db.mu.Lock()
	defer db.mu.Unlock()

	s := session.New(db.store, db.catalog)
	s.OnClosed(db.onSessionClosed)
	db.mu.sessions[s.ID()] = s
	return s
}

// Close destroys the  GraphEngine database instances and all sessions will be terminated.
func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, s := range db.mu.sessions {
		s.OnClosed(db.onSessionClosedLocked)
		s.Close()
	}

	return nil
}

func (db *DB) onSessionClosed(s *session.Session) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.onSessionClosedLocked(s)
}

func (db *DB) onSessionClosedLocked(s *session.Session) {
	delete(db.mu.sessions, s.ID())
}
