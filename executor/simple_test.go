// ---

package executor_test

import (
	"context"
	"testing"

	"github.com/simbiont-runtime/graphengine"
	"github.com/simbiont-runtime/graphengine/compiler"
	"github.com/simbiont-runtime/graphengine/parser"
	"github.com/simbiont-runtime/graphengine/session"
	"github.com/stretchr/testify/assert"
)

func TestSimpleExec(t *testing.T) {
	assert := assert.New(t)
	db, err := graphengine.Open(t.TempDir(), nil)
	assert.Nil(err)

	catalog := db.Catalog()

	cases := []struct {
		query string
		check func(session *session.Session)
	}{
		{
			query: "create graph g1",
			check: func(_ *session.Session) {
				assert.NotNil(catalog.Graph("g1"))
			},
		},
		{
			query: "use g1",
			check: func(s *session.Session) {
				assert.Equal("g1", s.StmtContext().CurrentGraphName())
			},
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		parser := parser.New()
		stmt, err := parser.ParseOneStmt(c.query)
		assert.Nil(err)

		s := db.NewSession()
		sc := s.StmtContext()
		exec, err := compiler.Compile(sc, stmt)
		assert.Nil(err)

		err = exec.Open(ctx)
		assert.Nil(err)
		_, err = exec.Next(ctx)
		assert.Nil(err)

		if c.check != nil {
			c.check(s)
		}
	}
}
