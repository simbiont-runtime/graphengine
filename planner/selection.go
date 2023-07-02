// ---

package planner

import "github.com/simbiont-runtime/graphengine/expression"

type LogicalSelection struct {
	baseLogicalPlan

	Condition expression.Expression
}

type PhysicalSelection struct {
	basePhysicalPlan

	Condition expression.Expression
}
