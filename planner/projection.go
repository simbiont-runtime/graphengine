// ---

package planner

import "github.com/simbiont-runtime/graphengine/expression"

type LogicalProjection struct {
	baseLogicalPlan

	Exprs []expression.Expression
}

type PhysicalProjection struct {
	basePhysicalPlan

	Exprs []expression.Expression
}
