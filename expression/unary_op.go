//  Copyright 2023  GraphEngine Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"fmt"

	"github.com/cockroachdb/apd/v3"
	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/parser/opcode"
	"github.com/simbiont-runtime/graphengine/stmtctx"
	"github.com/simbiont-runtime/graphengine/types"
)

var _ Expression = &UnaryExpr{}

type UnaryExpr struct {
	Op     opcode.Op
	Expr   Expression
	EvalOp UnaryEvalOp
}

func (u *UnaryExpr) String() string {
	return fmt.Sprintf("%s %s", u.Op, u.Expr)
}

func (u *UnaryExpr) ReturnType() types.T {
	return u.EvalOp.InferReturnType(u.Expr.ReturnType())
}

func (u *UnaryExpr) Eval(stmtCtx *stmtctx.Context, input datum.Row) (datum.Datum, error) {
	d, err := u.Expr.Eval(stmtCtx, input)
	if err != nil {
		return d, err
	}
	if d == datum.Null && !u.EvalOp.CallOnNullInput() {
		return d, nil
	}
	return u.EvalOp.Eval(stmtCtx, d)
}

func NewUnaryExpr(op opcode.Op, expr Expression) (*UnaryExpr, error) {
	unaryOp, ok := unaryOps[op]
	if !ok {
		return nil, fmt.Errorf("unsupported unary operator: %s", op)
	}
	return &UnaryExpr{
		Op:     op,
		Expr:   expr,
		EvalOp: unaryOp,
	}, nil
}

type UnaryEvalOp interface {
	InferReturnType(inputType types.T) types.T
	CallOnNullInput() bool
	Eval(stmtCtx *stmtctx.Context, input datum.Datum) (datum.Datum, error)
}

var unaryOps = map[opcode.Op]UnaryEvalOp{
	opcode.Minus: unaryMinusOp{},
}

type unaryMinusOp struct{}

func (u unaryMinusOp) InferReturnType(inputType types.T) types.T {
	return inputType
}

func (u unaryMinusOp) CallOnNullInput() bool {
	return false
}

func (u unaryMinusOp) Eval(_ *stmtctx.Context, input datum.Datum) (datum.Datum, error) {
	switch input.Type() {
	case types.Int:
		return datum.NewInt(-datum.AsInt(input)), nil
	case types.Float:
		return datum.NewFloat(-datum.AsFloat(input)), nil
	case types.Decimal:
		d := datum.AsDecimal(input)
		res := (&apd.Decimal{}).Neg(d)
		return datum.NewDecimal(res), nil
	default:
		return nil, fmt.Errorf("cannot negate %s", input.Type())
	}
}
