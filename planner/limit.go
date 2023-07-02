// ---

package planner

import "github.com/simbiont-runtime/graphengine/expression"

type LogicalLimit struct {
	baseLogicalPlan

	Offset expression.Expression
	Count  expression.Expression
}

type PhysicalLimit struct {
	basePhysicalPlan

	Offset expression.Expression
	Count  expression.Expression
}
