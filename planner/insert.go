// ---

package planner

import (
	"github.com/simbiont-runtime/graphengine/catalog"
	"github.com/simbiont-runtime/graphengine/expression"
	"github.com/simbiont-runtime/graphengine/parser/ast"
)

// Insert represents the plan of INSERT statement.
type Insert struct {
	basePlan

	Graph      *catalog.Graph
	Insertions []*ElementInsertion
	MatchPlan  Plan
}

// ElementInsertion represents a graph insertion element.
type ElementInsertion struct {
	Type ast.InsertionType
	// INSERT EDGE e BETWEEN x AND y FROM MATCH (x) , MATCH (y) WHERE id(x) = 1 AND id(y) = 2
	// FromIDExpr and ToIDExpr are the expressions to get the ID of the source and destination vertex.
	FromIDExpr  expression.Expression
	ToIDExpr    expression.Expression
	Labels      []*catalog.Label
	Assignments []*expression.Assignment
}
