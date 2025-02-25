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

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/simbiont-runtime/graphengine/stmtctx"
	"github.com/simbiont-runtime/graphengine/types"
)

var _ Expression = &PropertyAccess{}

// PropertyAccess represents a property access expression.
type PropertyAccess struct {
	Expr         Expression
	VariableName model.CIStr
	PropertyName model.CIStr
}

func (p *PropertyAccess) String() string {
	return p.VariableName.O + "." + p.PropertyName.O
}

func (p *PropertyAccess) ReturnType() types.T {
	// Property types are unknown until runtime.
	return types.Unknown
}

func (p *PropertyAccess) Eval(stmtCtx *stmtctx.Context, input datum.Row) (datum.Datum, error) {
	d, err := p.Expr.Eval(stmtCtx, input)
	if err != nil || d == datum.Null {
		return d, err
	}

	v, ok := d.(*datum.Vertex)
	if ok {
		d, ok := v.Props[p.PropertyName.L]
		if !ok {
			return datum.Null, nil
		}
		return d, nil
	}

	e, ok := d.(*datum.Edge)
	if ok {
		d, ok := e.Props[p.PropertyName.L]
		if !ok {
			return datum.Null, nil
		}
		return d, nil
	}

	return nil, fmt.Errorf("cannot access property on non-vertex or non-edge type %T", d)
}
