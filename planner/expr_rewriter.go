// ---

package planner

import (
	"fmt"

	"github.com/simbiont-runtime/graphengine/expression"
	"github.com/simbiont-runtime/graphengine/parser/ast"
)

type exprRewriter struct {
	p        LogicalPlan
	ctxStack []expression.Expression
	err      error
}

func RewriteExpr(expr ast.ExprNode, p LogicalPlan) (expression.Expression, error) {
	rewriter := &exprRewriter{
		p: p,
	}
	expr.Accept(rewriter)
	if rewriter.err != nil {
		return nil, rewriter.err
	}
	return rewriter.ctxStack[0], nil
}

// Enter implements the ast.Visitor interface.
func (er *exprRewriter) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	switch expr := n.(type) {
	case *ast.ValueExpr:
		_ = expr
	}
	return n, false
}

// Leave implements the ast.Visitor interface.
func (er *exprRewriter) Leave(n ast.Node) (node ast.Node, ok bool) {
	if er.err != nil {
		return node, false
	}

	switch expr := n.(type) {
	case *ast.ValueExpr:
		er.ctxStackAppend(&expression.Constant{Value: expr.Datum})
	case *ast.BinaryExpr:
		lExpr := er.ctxStack[er.ctxStackLen()-2]
		rExpr := er.ctxStack[er.ctxStackLen()-1]
		er.ctxStackPop(2)
		binExpr, err := expression.NewBinaryExpr(expr.Op, lExpr, rExpr)
		if err != nil {
			er.err = err
			return n, true
		}
		er.ctxStackAppend(binExpr)
	case *ast.UnaryExpr:
		input := er.ctxStack[er.ctxStackLen()-1]
		er.ctxStackPop(1)
		unaryExpr, err := expression.NewUnaryExpr(expr.Op, input)
		if err != nil {
			er.err = err
			return n, true
		}
		er.ctxStackAppend(unaryExpr)
	case *ast.VariableReference:
		idx := er.p.Columns().FindColumnIndex(expr.VariableName)
		if idx == -1 {
			er.err = fmt.Errorf("unresolved variable %s", expr.VariableName)
			return n, true
		}
		er.ctxStackAppend(&expression.Column{
			Index: idx,
			Name:  expr.VariableName,
			Type:  er.p.Columns()[idx].Type,
		})
	case *ast.PropertyAccess:
		idx := er.p.Columns().FindColumnIndex(expr.VariableName)
		if idx == -1 {
			er.err = fmt.Errorf("unresolved variable %s", expr.VariableName)
			return n, true
		}
		col := &expression.Column{
			Index: idx,
			Name:  expr.VariableName,
			Type:  er.p.Columns()[idx].Type,
		}
		er.ctxStackAppend(&expression.PropertyAccess{
			Expr:         col,
			VariableName: expr.VariableName,
			PropertyName: expr.PropertyName,
		})
	}

	return n, true
}

func (er *exprRewriter) ctxStackLen() int {
	return len(er.ctxStack)
}

func (er *exprRewriter) ctxStackPop(num int) {
	l := er.ctxStackLen()
	er.ctxStack = er.ctxStack[:l-num]
}

func (er *exprRewriter) ctxStackAppend(col expression.Expression) {
	er.ctxStack = append(er.ctxStack, col)
}
