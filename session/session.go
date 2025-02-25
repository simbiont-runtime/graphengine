// ---

package session

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/pingcap/errors"
	"github.com/simbiont-runtime/graphengine/catalog"
	"github.com/simbiont-runtime/graphengine/compiler"
	"github.com/simbiont-runtime/graphengine/parser"
	"github.com/simbiont-runtime/graphengine/parser/ast"
	"github.com/simbiont-runtime/graphengine/stmtctx"
	"github.com/simbiont-runtime/graphengine/storage/kv"
)

var (
	idGenerator atomic.Int64
	parserPool  = &sync.Pool{New: func() interface{} { return parser.New() }}
)

//	Session represents the session to interact with  GraphEngine database instance.
//
// Typically, the number of session will be same as the concurrent thread
// count of the application.
// All execution intermediate variables should be placed in the Context.
type Session struct {
	// Protect the current session will not be used concurrently.
	mu       sync.Mutex
	id       int64
	sc       *stmtctx.Context
	wg       sync.WaitGroup
	store    kv.Storage
	catalog  *catalog.Catalog
	closed   atomic.Bool
	cancelFn context.CancelFunc

	// Callback function while session closing.
	closeCallback func(s *Session)
}

// New returns a new session instance.
func New(store kv.Storage, catalog *catalog.Catalog) *Session {
	return &Session{
		id:      idGenerator.Add(1),
		sc:      stmtctx.New(store, catalog),
		store:   store,
		catalog: catalog,
	}
}

// ID returns a integer identifier of the current session.
func (s *Session) ID() int64 {
	return s.id
}

// StmtContext returns the statement context object.
func (s *Session) StmtContext() *stmtctx.Context {
	return s.sc
}

// Execute executes a query and reports whether the query executed successfully or not.
// A result set will be non-empty if execute successfully.
func (s *Session) Execute(ctx context.Context, query string) (ResultSet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancelFn := context.WithCancel(ctx)
	s.cancelFn = cancelFn
	s.wg.Add(1)
	defer s.wg.Done()

	p := parserPool.Get().(*parser.Parser)
	defer parserPool.Put(p)

	stmts, warns, err := p.Parse(query)
	if err != nil {
		return nil, err
	}
	for _, warn := range warns {
		s.sc.AppendWarning(errors.Annotate(warn, "parse warning"))
	}
	if len(stmts) == 0 {
		return emptyResultSet{}, nil
	}
	if len(stmts) > 1 {
		return nil, ErrMultipleStatementsNotSupported
	}

	return s.executeStmt(ctx, stmts[0])
}

func (s *Session) executeStmt(ctx context.Context, node ast.StmtNode) (ResultSet, error) {
	// TODO: support transaction

	// Reset the current statement context and prepare for executing the next statement.
	s.sc.Reset()

	exec, err := compiler.Compile(s.sc, node)
	if err != nil {
		return nil, err
	}
	err = exec.Open(ctx)
	if err != nil {
		return nil, err
	}

	return newQueryResultSet(exec), nil
}

// Close terminates the current session.
func (s *Session) Close() {
	if s.closed.Swap(true) {
		return
	}

	// Wait the current execution finished.
	if s.cancelFn != nil {
		s.cancelFn()
	}
	s.wg.Wait()

	if s.closeCallback != nil {
		s.closeCallback(s)
	}
}

// OnClosed sets the closed callback which will invoke after session closed.
func (s *Session) OnClosed(cb func(session *Session)) {
	s.mu.Lock()
	s.mu.Unlock()
	s.closeCallback = cb
}
