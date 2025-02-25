// ---

package ast

import (
	"github.com/simbiont-runtime/graphengine/parser/format"
	"github.com/simbiont-runtime/graphengine/parser/model"
)

var (
	_ Node = &UseStmt{}
	_ Node = &BeginStmt{}
	_ Node = &RollbackStmt{}
	_ Node = &CommitStmt{}
	_ Node = &ExplainStmt{}
	_ Node = &ShowStmt{}
)

type UseStmt struct {
	stmtNode

	GraphName model.CIStr
}

func (u *UseStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("USE ")
	ctx.WriteName(u.GraphName.String())
	return nil
}

func (u *UseStmt) Accept(v Visitor) (node Node, ok bool) {
	newNode, skipChildren := v.Enter(u)
	if skipChildren {
		return v.Leave(newNode)
	}
	return v.Leave(newNode)
}

type BeginStmt struct {
	stmtNode
}

func (b *BeginStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("BEGIN")
	return nil
}

func (b *BeginStmt) Accept(v Visitor) (node Node, ok bool) {
	newNode, skipChildren := v.Enter(b)
	if skipChildren {
		return v.Leave(newNode)
	}
	return v.Leave(newNode)
}

type RollbackStmt struct {
	node
}

func (r *RollbackStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("ROLLBACK")
	return nil
}

func (r *RollbackStmt) Accept(v Visitor) (node Node, ok bool) {
	newNode, skipChildren := v.Enter(r)
	if skipChildren {
		return v.Leave(newNode)
	}
	return v.Leave(newNode)
}

type CommitStmt struct {
	stmtNode
}

func (c *CommitStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("COMMIT")
	return nil
}

func (c *CommitStmt) Accept(v Visitor) (node Node, ok bool) {
	newNode, skipChildren := v.Enter(c)
	if skipChildren {
		return v.Leave(newNode)
	}
	return v.Leave(newNode)
}

type ExplainStmt struct {
	stmtNode

	Select *SelectStmt
}

func (e *ExplainStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("EXPLAIN ")
	return e.Select.Restore(ctx)
}

func (e *ExplainStmt) Accept(v Visitor) (node Node, ok bool) {
	newNode, skipChildren := v.Enter(e)
	if skipChildren {
		return v.Leave(newNode)
	}

	nn := newNode.(*ExplainStmt)
	n, ok := nn.Select.Accept(v)
	if !ok {
		return nn, false
	}
	nn.Select = n.(*SelectStmt)

	return v.Leave(nn)
}

type ShowTarget byte

const (
	ShowTargetGraphs ShowTarget = iota + 1
	ShowTargetLabels
)

type ShowStmt struct {
	stmtNode

	Tp        ShowTarget
	GraphName model.CIStr
}

func (s *ShowStmt) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("SHOW ")
	switch s.Tp {
	case ShowTargetGraphs:
		ctx.WriteKeyWord("GRAPHS")
	case ShowTargetLabels:
		ctx.WriteKeyWord("LABELS")
		if !s.GraphName.IsEmpty() {
			ctx.WriteKeyWord(" IN ")
			ctx.WriteName(s.GraphName.String())
		}
	}
	return nil
}

func (s *ShowStmt) Accept(v Visitor) (node Node, ok bool) {
	newNode, skipChildren := v.Enter(s)
	if skipChildren {
		return v.Leave(newNode)
	}
	return v.Leave(newNode)
}
