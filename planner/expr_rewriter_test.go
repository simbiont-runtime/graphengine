// ---

package planner_test

import (
	"testing"

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/expression"
	"github.com/simbiont-runtime/graphengine/parser/ast"
	"github.com/simbiont-runtime/graphengine/planner"
	"github.com/stretchr/testify/assert"
)

func TestRewriteExpr(t *testing.T) {
	cases := []struct {
		expr   ast.ExprNode
		expect expression.Expression
	}{
		{
			expr:   &ast.ValueExpr{Datum: datum.NewInt(1)},
			expect: &expression.Constant{Value: datum.NewInt(1)},
		},
	}

	for _, c := range cases {
		expr, err := planner.RewriteExpr(c.expr, &planner.LogicalDual{})
		assert.Nil(t, err)
		assert.Equal(t, c.expect, expr)
	}
}
