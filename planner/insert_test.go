// ---

package planner_test

import (
	"sync/atomic"
	"testing"

	"github.com/simbiont-runtime/graphengine"
	"github.com/simbiont-runtime/graphengine/catalog"
	"github.com/simbiont-runtime/graphengine/parser"
	"github.com/simbiont-runtime/graphengine/parser/ast"
	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/simbiont-runtime/graphengine/planner"
	"github.com/simbiont-runtime/graphengine/stmtctx"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildInsert(t *testing.T) {
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
		check func(insert *planner.Insert)
	}{
		// Catalog information refer: initCatalog
		{
			query: "insert vertex x labels(A, B, C)",
			check: func(insert *planner.Insert) {
				assert.Equal(1, len(insert.Insertions))
				assert.Equal(3, len(insert.Insertions[0].Labels))
			},
		},
		{
			query: "insert into graph102 vertex x labels(A, B, C)",
			check: func(insert *planner.Insert) {
				assert.Equal("graph102", insert.Graph.Meta().Name.L)
			},
		},
		{
			query: "insert vertex x properties (x.name = 'test')",
			check: func(insert *planner.Insert) {
				assert.Equal(1, len(insert.Insertions[0].Assignments))
			},
		},
		{
			query: `insert vertex x labels(A, B, C) properties (x.name = 'test'),
					       vertex y labels(A, B) properties (y.name = 'test'),
					       vertex z labels(B, C) properties (z.name = 'test')`,
			check: func(insert *planner.Insert) {
				assert.Equal(3, len(insert.Insertions))
			},
		},
		{
			query: `insert vertex x labels(A, B, C) properties (x.name = 'test'),
					       vertex y labels(A, B) properties (y.name = 'test'),
					       edge z between x and y labels(B, C) from match (x), match (y)`,
			check: func(insert *planner.Insert) {
				assert.Equal(3, len(insert.Insertions))
				assert.Equal(ast.InsertionTypeEdge, insert.Insertions[2].Type)
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
		insert, ok := plan.(*planner.Insert)
		assert.True(ok)
		c.check(insert)
	}
}
