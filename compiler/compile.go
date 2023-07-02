// ---

package compiler

import (
	"github.com/simbiont-runtime/graphengine/executor"
	"github.com/simbiont-runtime/graphengine/parser/ast"
	"github.com/simbiont-runtime/graphengine/planner"
	"github.com/simbiont-runtime/graphengine/stmtctx"
)

// Compile compiles the statement AST node into an executable statement. The compiler relay
// on the statement context to retrieve some environment information and set some intermediate
// variables while compiling. The catalog is used to resolve names in the query.
func Compile(sc *stmtctx.Context, node ast.StmtNode) (executor.Executor, error) {
	// Macro expansion
	macroExp := NewMacroExpansion()
	node.Accept(macroExp)

	// Check the AST to ensure it is valid.
	preprocess := NewPreprocess(sc)
	node.Accept(preprocess)
	err := preprocess.Error()
	if err != nil {
		return nil, err
	}

	// Prepare missing properties
	propPrep := NewPropertyPreparation(sc)
	node.Accept(propPrep)
	err = propPrep.CreateMissing()
	if err != nil {
		return nil, err
	}

	// Build plan tree from a valid AST.
	planBuilder := planner.NewBuilder(sc)
	plan, err := planBuilder.Build(node)
	if err != nil {
		return nil, err
	}
	logicalPlan, isLogicalPlan := plan.(planner.LogicalPlan)
	if isLogicalPlan {
		// Optimize the logical plan and generate physical plan.
		plan = planner.Optimize(logicalPlan)
	}

	execBuilder := executor.NewBuilder(sc)
	exec := execBuilder.Build(plan)
	err = execBuilder.Error()
	if err != nil {
		return nil, err
	}

	return exec, nil
}
