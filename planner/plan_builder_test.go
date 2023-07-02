// ---

package planner_test

import (
	"testing"

	"github.com/simbiont-runtime/graphengine"
	"github.com/simbiont-runtime/graphengine/parser"
	"github.com/simbiont-runtime/graphengine/planner"
	"github.com/simbiont-runtime/graphengine/stmtctx"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildDDL(t *testing.T) {
	assert := assert.New(t)

	cases := []string{
		"create graph if not exists graph5",
		"create graph graph5",
		"create label label1",
		"create label if not exists label1",
		"create index index1 (a, b)",
		"create index if not exists index1 (a, b)",
		"drop graph graph5",
		"drop label label1",
		"drop index index1",
	}

	db, err := graphengine.Open(t.TempDir(), nil)
	assert.Nil(err)

	for _, c := range cases {
		parser := parser.New()
		stmt, err := parser.ParseOneStmt(c)
		assert.Nil(err)

		builder := planner.NewBuilder(stmtctx.New(db.Store(), db.Catalog()))
		plan, err := builder.Build(stmt)
		assert.Nil(err)

		ddl, ok := plan.(*planner.DDL)
		assert.True(ok)
		assert.Equal(stmt, ddl.Statement)
	}
}

func TestBuilder_BuildSimple(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		query string
	}{
		// Catalog information refer: initCatalog
		{
			query: "use graph100",
		},
		{
			query: "use graph1",
		},
	}

	db, err := graphengine.Open(t.TempDir(), nil)
	assert.Nil(err)

	for _, c := range cases {
		parser := parser.New()
		stmt, err := parser.ParseOneStmt(c.query)
		assert.Nil(err)

		builder := planner.NewBuilder(stmtctx.New(db.Store(), db.Catalog()))
		plan, err := builder.Build(stmt)
		assert.Nil(err)
		ddl, ok := plan.(*planner.Simple)
		assert.True(ok)
		assert.Equal(stmt, ddl.Statement)
	}
}
