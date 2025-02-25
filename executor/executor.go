// ---

package executor

import (
	"context"

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/planner"
	"github.com/simbiont-runtime/graphengine/stmtctx"
)

// Executor is the physical implementation of an algebra operator.
type Executor interface {
	base() *baseExecutor
	Columns() planner.ResultColumns
	Open(context.Context) error
	Next(context.Context) (datum.Row, error)
	Close() error
}

type baseExecutor struct {
	sc       *stmtctx.Context
	id       int
	columns  planner.ResultColumns // output columns
	children []Executor
}

func newBaseExecutor(sc *stmtctx.Context, cols planner.ResultColumns, id int, children ...Executor) baseExecutor {
	e := baseExecutor{
		children: children,
		sc:       sc,
		id:       id,
		columns:  cols,
	}
	return e
}

// base returns the baseExecutor of an executor, don't override this method!
func (e *baseExecutor) base() *baseExecutor {
	return e
}

// Open initializes children recursively and "childrenResults" according to children's schemas.
func (e *baseExecutor) Open(ctx context.Context) error {
	for _, child := range e.children {
		err := child.Open(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Next fills multiple rows into a chunk.
func (e *baseExecutor) Next(context.Context) (datum.Row, error) {
	return nil, nil
}

// Close closes all executors and release all resources.
func (e *baseExecutor) Close() error {
	var firstErr error
	for _, src := range e.children {
		if err := src.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (e *baseExecutor) Columns() planner.ResultColumns {
	return e.columns
}
