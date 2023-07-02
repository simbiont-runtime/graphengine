// ---

//go:build ignore

package planner_test

import (
	"sync/atomic"
	"testing"

	"github.com/simbiont-runtime/graphengine"
	"github.com/simbiont-runtime/graphengine/catalog"
	"github.com/simbiont-runtime/graphengine/parser"
	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/simbiont-runtime/graphengine/planner"
	"github.com/simbiont-runtime/graphengine/stmtctx"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildMatch(t *testing.T) {
	assert := assert.New(t)

	db, err := graphengine.Open(t.TempDir(), nil)
	assert.Nil(err)

	// Prepare mock catalog.
	mockGraphs := []string{"graph101", "graph102"}
	mockLabels := []string{"A", "B", "C"}
	mockProps := []string{"name", "age"}
	id := atomic.Int64{}
	for _, g := range mockGraphs {
		graphID := id.Add(1)
		db.Catalog().Apply(&catalog.Patch{
			Type: catalog.PatchTypeCreateGraph,
			Data: &model.GraphInfo{
				ID:   graphID,
				Name: model.NewCIStr(g),
			},
		})
		for _, l := range mockLabels {
			db.Catalog().Apply(&catalog.Patch{
				Type: catalog.PatchTypeCreateLabel,
				Data: &catalog.PatchLabel{
					GraphID: graphID,
					LabelInfo: &model.LabelInfo{
						ID:   id.Add(1),
						Name: model.NewCIStr(l),
					},
				},
			})

		}
		var properties []*model.PropertyInfo
		for i, p := range mockProps {
			properties = append(properties, &model.PropertyInfo{
				ID:   uint16(i + 1),
				Name: model.NewCIStr(p),
			})
		}
		db.Catalog().Apply(&catalog.Patch{
			Type: catalog.PatchTypeCreateProperties,
			Data: &catalog.PatchProperties{
				MaxPropID:  uint16(len(properties)),
				GraphID:    graphID,
				Properties: properties,
			},
		})

	}

	cases := []struct {
		query string
		check func(proj *planner.LogicalProjection)
	}{
		// Catalog information refer: initCatalog
		{
			query: "SELECT * FROM MATCH (n:A)->(m:B)",
			check: func(proj *planner.LogicalProjection) {
				match := proj.Children()[0].(*planner.LogicalMatch)
				assert.Equal(1, len(match.Subgraphs))
				assert.Equal(1, len(match.Subgraphs[0].Paths))
				assert.Equal(2, len(match.Subgraphs[0].Paths[0].Vertices))
			},
		},
		{
			query: "SELECT * FROM MATCH ( (n:A)->(m:B), (c:C)->(m:B) )",
			check: func(proj *planner.LogicalProjection) {
				match := proj.Children()[0].(*planner.LogicalMatch)
				assert.Equal(1, len(match.Subgraphs))
				assert.Equal(2, len(match.Subgraphs[0].Paths))
				assert.Equal(2, len(match.Subgraphs[0].Paths[0].Vertices))
			},
		},
		{
			query: "SELECT * FROM MATCH (n:A)->(m:B), MATCH (c:C)->(m:B)",
			check: func(proj *planner.LogicalProjection) {
				match := proj.Children()[0].(*planner.LogicalMatch)
				assert.Equal(2, len(match.Subgraphs))
				assert.Equal(1, len(match.Subgraphs[0].Paths))
				assert.Equal(2, len(match.Subgraphs[0].Paths[0].Vertices))
			},
		},
	}

	for _, c := range cases {
		parser := parser.New()
		stmt, err := parser.ParseOneStmt(c.query)
		assert.Nil(err, c.query)

		sc := stmtctx.New(db.Store(), db.Catalog())
		sc.SetCurrentGraphName("graph101")

		builder := planner.NewBuilder(sc)
		plan, err := builder.Build(stmt)
		assert.Nil(err)
		projection, ok := plan.(*planner.LogicalProjection)
		assert.True(ok)
		c.check(projection)
	}
}
