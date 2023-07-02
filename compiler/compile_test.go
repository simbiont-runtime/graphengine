// ---

package compiler_test

import (
	"testing"

	"github.com/simbiont-runtime/graphengine"
	"github.com/simbiont-runtime/graphengine/catalog"
	"github.com/simbiont-runtime/graphengine/compiler"
	"github.com/simbiont-runtime/graphengine/executor"
	"github.com/simbiont-runtime/graphengine/parser"
	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	t.Skip()
	assert := assert.New(t)
	db, err := graphengine.Open(t.TempDir(), nil)
	assert.Nil(err)

	db.Catalog().Apply(&catalog.Patch{
		Type: catalog.PatchTypeCreateGraph,
		Data: &model.GraphInfo{
			ID:   1,
			Name: model.NewCIStr("g1"),
		},
	})

	ddl := func(exec executor.Executor) {
		_, ok := exec.(*executor.DDLExec)
		assert.True(ok)
	}

	simple := func(exec executor.Executor) {
		_, ok := exec.(*executor.SimpleExec)
		assert.True(ok)
	}

	cases := []struct {
		query string
		check func(exec executor.Executor)
	}{
		{
			query: "create graph g2",
			check: ddl,
		},
		{
			query: "create graph if not exists g1",
			check: ddl,
		},
		{
			query: "create label if not exists l1",
			check: ddl,
		},
		{
			query: "create index if not exists i1 (a)",
			check: ddl,
		},
		{
			query: "use g1",
			check: simple,
		},
	}

	sc := db.NewSession().StmtContext()
	sc.SetCurrentGraphName("g1")

	for _, c := range cases {
		parser := parser.New()
		stmt, err := parser.ParseOneStmt(c.query)
		assert.Nil(err, c.query)
		exec, err := compiler.Compile(sc, stmt)
		assert.Nil(err, c.query)
		c.check(exec)
	}
}
