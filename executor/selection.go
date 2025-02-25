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

package executor

import (
	"context"

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/expression"
)

// SelectionExec represents a selection executor.
type SelectionExec struct {
	baseExecutor

	condition expression.Expression
}

func (p *SelectionExec) Next(ctx context.Context) (datum.Row, error) {
	for {
		row, err := p.children[0].Next(ctx)
		if err != nil || row == nil {
			return nil, err
		}
		d, err := p.condition.Eval(p.sc, row)
		if err != nil {
			return nil, err
		}
		if datum.AsBool(d) {
			return row, nil
		}
	}
}
