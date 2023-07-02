// ---

package planner

// Optimize optimizes the plan to the optimal physical plan.
func Optimize(plan LogicalPlan) Plan {
	switch p := plan.(type) {
	case *LogicalMatch:
		return optimizeMatch(p)
	case *LogicalProjection:
		return optimizeProjection(p)
	case *LogicalSelection:
		return optimizeSelection(p)
	}
	return plan
}

func optimizeMatch(plan *LogicalMatch) Plan {
	result := &PhysicalMatch{}
	result.SetColumns(plan.Columns())
	result.Subgraph = plan.Subgraph
	return result
}

func optimizeProjection(plan *LogicalProjection) Plan {
	result := &PhysicalProjection{}
	result.SetColumns(plan.Columns())
	result.Exprs = plan.Exprs
	childPlan := Optimize(plan.Children()[0])
	result.SetChildren(childPlan.(PhysicalPlan))
	return result
}

func optimizeSelection(plan *LogicalSelection) Plan {
	result := &PhysicalSelection{}
	result.SetColumns(plan.Columns())
	result.Condition = plan.Condition
	childPlan := Optimize(plan.Children()[0])
	result.SetChildren(childPlan.(PhysicalPlan))
	return result
}
