// ---

package compiler

import (
	"testing"

	"github.com/simbiont-runtime/graphengine/parser"
	"github.com/simbiont-runtime/graphengine/parser/ast"
	"github.com/simbiont-runtime/graphengine/parser/opcode"
	"github.com/stretchr/testify/assert"
)

func TestMacroExpansion(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		query string
		check func(node ast.Node)
	}{
		{
			query: `PATH has_parent AS () -[:has_father|has_mother]-> (:Person)
					SELECT ancestor.name
					  FROM MATCH (p1:Person) -/:has_parent+/-> (ancestor)
						 , MATCH (p2:Person) -/:has_parent+/-> (ancestor)
					 WHERE p1.name = 'Mario'
					   AND p2.name = 'Luigi'`,
			check: func(node ast.Node) {
				stmt := node.(*ast.SelectStmt)
				assert.Equal(1, len(stmt.PathPatternMacros))
				expr := stmt.From.Matches[0].Paths[0].Connections[0].(*ast.ReachabilityPathExpr)
				assert.Equal(expr.Macros["has_parent"], stmt.PathPatternMacros[0].Path)
			},
		},
		{
			query: `PATH connects_to AS (:Generator) -[:has_connector]-> (c:Connector) <-[:has_connector]- (:Generator)
					  WHERE c.status = 'OPERATIONAL'
					SELECT generatorA.location, generatorB.location
					  FROM MATCH (generatorA) -/:connects_to+/-> (generatorB)`,
			check: func(node ast.Node) {
				stmt := node.(*ast.SelectStmt)
				assert.Equal(1, len(stmt.PathPatternMacros))
				expr := stmt.From.Matches[0].Paths[0].Connections[0].(*ast.ReachabilityPathExpr)
				assert.Equal(expr.Macros["connects_to"], stmt.PathPatternMacros[0].Path)
				assert.NotNil(stmt.Where)
			},
		},
		{
			query: `PATH connects_to AS (:Generator) -[:has_connector]-> (c:Connector) <-[:has_connector]- (:Generator)
					  WHERE c.status = 'OPERATIONAL'
					PATH has_parent AS () -[:has_father|has_mother]-> (:Person)
					SELECT generatorA.location, generatorB.location
					  FROM MATCH (generatorA) -/:connects_to|has_parent+/-> (generatorB)`,
			check: func(node ast.Node) {
				stmt := node.(*ast.SelectStmt)
				assert.Equal(2, len(stmt.PathPatternMacros))
				expr := stmt.From.Matches[0].Paths[0].Connections[0].(*ast.ReachabilityPathExpr)
				assert.Equal(expr.Macros["connects_to"], stmt.PathPatternMacros[0].Path)
				assert.Equal(expr.Macros["has_parent"], stmt.PathPatternMacros[1].Path)
				assert.NotNil(stmt.Where)
			},
		},
		{
			query: `PATH connects_to AS (:Generator) -[:has_connector]-> (c:Connector) <-[:has_connector]- (:Generator)
					  WHERE c.status = 'OPERATIONAL'
					PATH has_parent AS () -[f:has_father|has_mother]-> (:Person)
					  WHERE f.age > 30
					SELECT generatorA.location, generatorB.location
					  FROM MATCH (generatorA) -/:connects_to|has_parent+/-> (generatorB)`,
			check: func(node ast.Node) {
				stmt := node.(*ast.SelectStmt)
				assert.Equal(2, len(stmt.PathPatternMacros))
				expr := stmt.From.Matches[0].Paths[0].Connections[0].(*ast.ReachabilityPathExpr)
				assert.Equal(expr.Macros["connects_to"], stmt.PathPatternMacros[0].Path)
				assert.Equal(expr.Macros["has_parent"], stmt.PathPatternMacros[1].Path)
				assert.NotNil(stmt.Where)
				logicalAnd := stmt.Where.(*ast.BinaryExpr)
				assert.Equal(opcode.LogicAnd, logicalAnd.Op)
			},
		},
		{
			query: `PATH connects_to AS (:Generator) -[:has_connector]-> (c:Connector) <-[:has_connector]- (:Generator)
					  WHERE c.status = 'OPERATIONAL'
					PATH has_parent AS () -[f:has_father|has_mother]-> (:Person)
					  WHERE f.age > 30
					SELECT generatorA.location, generatorB.location
					  FROM MATCH (generatorA) -/:connects_to|has_parent+/-> (generatorB)
					 WHERE a > 10`,
			check: func(node ast.Node) {
				stmt := node.(*ast.SelectStmt)
				assert.Equal(2, len(stmt.PathPatternMacros))
				expr := stmt.From.Matches[0].Paths[0].Connections[0].(*ast.ReachabilityPathExpr)
				assert.Equal(expr.Macros["connects_to"], stmt.PathPatternMacros[0].Path)
				assert.Equal(expr.Macros["has_parent"], stmt.PathPatternMacros[1].Path)
				assert.NotNil(stmt.Where)
				logicalAnd := stmt.Where.(*ast.BinaryExpr)
				assert.Equal(opcode.LogicAnd, logicalAnd.Op)
				logicalEq := logicalAnd.L.(*ast.BinaryExpr)
				assert.Equal(opcode.GT, logicalEq.Op) //  a > 10
				logicalAnd = logicalAnd.R.(*ast.BinaryExpr)
				assert.Equal(opcode.LogicAnd, logicalAnd.Op)
			},
		},
	}

	for _, c := range cases {
		parser := parser.New()
		stmt, err := parser.ParseOneStmt(c.query)
		assert.Nil(err)

		exp := NewMacroExpansion()
		n, ok := stmt.Accept(exp)
		assert.True(ok)
		c.check(n)
	}
}
