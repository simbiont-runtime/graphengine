// ---

package graphengine

import (
	"context"
	"testing"

	"github.com/simbiont-runtime/graphengine/session"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	db, err := Open(t.TempDir(), nil)
	require.NoError(t, err)
	require.NotNil(t, db)
}

type TestKit struct {
	t    *testing.T
	sess *session.Session
}

func NewTestKit(t *testing.T, sess *session.Session) *TestKit {
	return &TestKit{
		t:    t,
		sess: sess,
	}
}

func (tk *TestKit) MustExec(ctx context.Context, query string) {
	rs, err := tk.sess.Execute(ctx, query)
	require.NoError(tk.t, err)
	require.NoError(tk.t, rs.Next(ctx))
}

func TestDDL(t *testing.T) {
	db, err := Open(t.TempDir(), nil)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	catalog := db.Catalog()
	sess := db.NewSession()
	require.NotNil(t, sess)

	tk := NewTestKit(t, sess)

	ctx := context.Background()
	tk.MustExec(ctx, "CREATE GRAPH graph101")
	require.NoError(t, err)
	graph := catalog.Graph("graph101")
	require.NotNil(t, graph)

	sess.StmtContext().SetCurrentGraphName("graph101")
	tk.MustExec(ctx, "CREATE LABEL label01")
	require.NoError(t, err)
	require.NotNil(t, graph.Label("label01"))

	tk.MustExec(ctx, "CREATE LABEL IF NOT EXISTS label01")
	require.NoError(t, err)

	tk.MustExec(ctx, "DROP LABEL label01")
	require.NoError(t, err)
	require.Nil(t, graph.Label("label01"))

	tk.MustExec(ctx, "DROP LABEL IF EXISTS label01")
	require.NoError(t, err)
	require.Nil(t, graph.Label("label01"))

	tk.MustExec(ctx, "DROP GRAPH graph101")
	require.NoError(t, err)
	require.Nil(t, catalog.Graph("graph101"))
}
