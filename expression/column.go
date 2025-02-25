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

var _ Expression = &Column{}

type Column struct {
	Index int
	Name  model.CIStr
	Type  types.T
}

func (c *Column) String() string {
	return c.Name.O
}

func (c *Column) ReturnType() types.T {
	return c.Type
}

func (c *Column) Eval(stmtCtx *stmtctx.Context, input datum.Row) (datum.Datum, error) {
	if c.Index >= len(input) {
		return nil, fmt.Errorf("column index %d out of input row length %d", c.Index, len(input))
	}
	return input[c.Index], nil
}
