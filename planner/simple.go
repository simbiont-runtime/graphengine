// ---

package planner

import "github.com/simbiont-runtime/graphengine/parser/ast"

// Simple represents the physical plan of simple statements.
type Simple struct {
	basePlan

	Statement ast.StmtNode
}
