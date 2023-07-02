// ---

package planner

import (
	"fmt"

	"github.com/simbiont-runtime/graphengine/expression"
	"github.com/simbiont-runtime/graphengine/parser/model"
)

// ByItem wraps a "by" item.
type ByItem struct {
	Expr      expression.Expression
	AsName    model.CIStr
	Desc      bool
	NullOrder bool
}

// String implements fmt.Stringer interface.
func (by *ByItem) String() string {
	if by.Desc {
		return fmt.Sprintf("%s true", by.Expr)
	}
	return by.Expr.String()
}

type LogicalSort struct {
	baseLogicalPlan

	ByItems []*ByItem
}

type PhysicalSort struct {
	basePhysicalPlan

	ByItems []*ByItem
}
