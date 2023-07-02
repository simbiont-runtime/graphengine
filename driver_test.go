package graphengine_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestDriver(t *testing.T) {
	db, err := sql.Open("graphEngine", t.TempDir())
	require.NoError(t, err)
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn := lo.Must1(db.Conn(ctx))
	_ = lo.Must1(conn.ExecContext(ctx, "CREATE GRAPH g"))
	_ = lo.Must1(conn.ExecContext(ctx, "USE g"))
	_ = lo.Must1(conn.ExecContext(ctx, "INSERT VERTEX x PROPERTIES (x.a = 123)"))
	rows := lo.Must1(conn.QueryContext(ctx, "SELECT x.a FROM MATCH (x)"))

	var a int
	require.True(t, rows.Next())
	require.NoError(t, rows.Scan(&a))
	require.Equal(t, 123, a)
	require.False(t, rows.Next())
	require.NoError(t, rows.Err())
}
