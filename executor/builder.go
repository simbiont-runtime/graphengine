// ---

package executor

import (
	"github.com/pingcap/errors"
	"github.com/simbiont-runtime/graphengine/codec"
	"github.com/simbiont-runtime/graphengine/parser/ast"
	"github.com/simbiont-runtime/graphengine/planner"
	"github.com/simbiont-runtime/graphengine/stmtctx"
)

// Builder is used to build from a plan into executor.
type Builder struct {
	sc  *stmtctx.Context
	err error
}

// NewBuilder returns a build instance.
func NewBuilder(sc *stmtctx.Context) *Builder {
	return &Builder{
		sc: sc,
	}
}

// Build builds an executor from a plan.
func (b *Builder) Build(plan planner.Plan) Executor {
	switch p := plan.(type) {
	case *planner.DDL:
		return b.buildDDL(p)
	case *planner.Simple:
		return b.buildSimple(p)
	case *planner.Insert:
		return b.buildInsert(p)
	case *planner.PhysicalMatch:
		return b.buildMatch(p)
	case *planner.PhysicalProjection:
		return b.buildProjection(p)
	case *planner.PhysicalSelection:
		return b.buildSelection(p)
	default:
		b.err = errors.Errorf("unknown plan: %T", plan)
	}
	return nil
}

// Error returns the internal error encountered while building.
func (b *Builder) Error() error {
	return b.err
}

func (b *Builder) buildDDL(plan *planner.DDL) Executor {
	exec := &DDLExec{
		baseExecutor: newBaseExecutor(b.sc, plan.Columns(), plan.ID()),
		statement:    plan.Statement,
	}
	return exec
}

func (b *Builder) buildSimple(plan *planner.Simple) Executor {
	var exec Executor
	switch s := plan.Statement.(type) {
	case *ast.ShowStmt:
		exec = &ShowExec{
			baseExecutor: newBaseExecutor(b.sc, showStmtColumns[s.Tp], plan.ID()),
			statement:    s,
		}
	default:
		exec = &SimpleExec{
			baseExecutor: newBaseExecutor(b.sc, plan.Columns(), plan.ID()),
			statement:    plan.Statement,
		}
	}
	return exec
}

func (b *Builder) buildInsert(plan *planner.Insert) Executor {
	exec := &InsertExec{
		baseExecutor: newBaseExecutor(b.sc, plan.Columns(), plan.ID()),
		graph:        plan.Graph,
		insertions:   plan.Insertions,
		encoder:      &codec.PropertyEncoder{},
		decoder:      &codec.PropertyDecoder{},
	}
	if plan.MatchPlan != nil {
		exec.matchExec = b.Build(plan.MatchPlan)
	}
	return exec
}

func (b *Builder) buildMatch(plan *planner.PhysicalMatch) Executor {
	exec := &MatchExec{
		baseExecutor: newBaseExecutor(b.sc, plan.Columns(), plan.ID()),
		subgraph:     plan.Subgraph,
	}
	return exec
}

func (b *Builder) buildProjection(plan *planner.PhysicalProjection) Executor {
	childExec := b.Build(plan.Children()[0])
	exec := &ProjectionExec{
		baseExecutor: newBaseExecutor(b.sc, plan.Columns(), plan.ID(), childExec),
		exprs:        plan.Exprs,
	}
	return exec
}

func (b *Builder) buildSelection(plan *planner.PhysicalSelection) Executor {
	childExec := b.Build(plan.Children()[0])
	exec := &SelectionExec{
		baseExecutor: newBaseExecutor(b.sc, plan.Columns(), plan.ID(), childExec),
		condition:    plan.Condition,
	}
	return exec
}
