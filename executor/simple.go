// ---

package executor

import (
	"context"

	"github.com/pingcap/errors"
	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/meta"
	"github.com/simbiont-runtime/graphengine/parser/ast"
)

// SimpleExec is used to execute some simple tasks.
type SimpleExec struct {
	baseExecutor

	done      bool
	statement ast.StmtNode
}

func (e *SimpleExec) Next(_ context.Context) (datum.Row, error) {
	if e.done {
		return nil, nil
	}
	e.done = true

	switch stmt := e.statement.(type) {
	case *ast.UseStmt:
		return nil, e.execUse(stmt)
	default:
		return nil, errors.Errorf("unknown statement: %T", e.statement)
	}
}

func (e *SimpleExec) execUse(stmt *ast.UseStmt) error {
	graph := e.sc.Catalog().Graph(stmt.GraphName.L)
	if graph == nil {
		return meta.ErrGraphNotExists
	}

	e.sc.SetCurrentGraphName(stmt.GraphName.L)

	return nil
}
