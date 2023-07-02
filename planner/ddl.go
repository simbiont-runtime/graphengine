// ---

package planner

import "github.com/simbiont-runtime/graphengine/parser/ast"

// DDL represents the physical plan of DDL statement.
type DDL struct {
	basePlan

	Statement ast.DDLNode
}
