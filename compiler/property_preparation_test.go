// ---

package compiler_test

import (
	"testing"

	"github.com/simbiont-runtime/graphengine"
	"github.com/simbiont-runtime/graphengine/compiler"
	"github.com/simbiont-runtime/graphengine/parser"
	"github.com/simbiont-runtime/graphengine/stmtctx"
	"github.com/stretchr/testify/assert"
)

func TestPropertyPreparation(t *testing.T) {
	assert := assert.New(t)
	tempDir := t.TempDir()
	initCatalog(assert, tempDir)

	db, err := graphengine.Open(tempDir, nil)
	assert.Nil(err)

	cases := []struct {
		query string
		check func()
	}{
		{
			query: "INSERT INTO graph1 VERTEX x LABELS (label1) PROPERTIES ( x.prop1 = 'test')",
			check: func() {
				graph := db.Catalog().Graph("graph1")
				assert.NotNil(graph.Property("prop1"))
			},
		},
		{
			query: "INSERT INTO graph2 VERTEX x LABELS (label1) PROPERTIES ( x.property = 'test')",
			check: func() {
				graph := db.Catalog().Graph("graph1")
				assert.NotNil(2, len(graph.Properties()))
			},
		},
		{
			query: "INSERT INTO graph2 VERTEX x LABELS (label1) PROPERTIES ( x.property2 = 'test')",
			check: func() {
				graph := db.Catalog().Graph("graph1")
				assert.NotNil(3, len(graph.Properties()))
				assert.NotNil(graph.Property("property2"))
			},
		},
	}

	for _, c := range cases {
		parser := parser.New()
		stmt, err := parser.ParseOneStmt(c.query)
		assert.Nil(err)
		sc := stmtctx.New(db.Store(), db.Catalog())

		prep := compiler.NewPropertyPreparation(sc)
		stmt.Accept(prep)
		err = prep.CreateMissing()
		assert.Nil(err)
	}
}
