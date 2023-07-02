// ---

package executor

import (
	"context"

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/expression"
)

// ProjectionExec represents a projection executor.
type ProjectionExec struct {
	baseExecutor

	exprs []expression.Expression
}

func (p *ProjectionExec) Next(ctx context.Context) (datum.Row, error) {
	childRow, err := p.children[0].Next(ctx)
	if err != nil || childRow == nil {
		return nil, err
	}

	result := make(datum.Row, len(p.exprs))
	for i, expr := range p.exprs {
		d, err := expr.Eval(p.sc, childRow)
		if err != nil {
			return nil, err
		}
		result[i] = d
	}
	return result, nil
}
