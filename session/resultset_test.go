// ---

package session

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyResultSet(t *testing.T) {
	rs := emptyResultSet{}
	require.NoError(t, rs.Next(context.Background()))
	require.False(t, rs.Valid())
	require.Empty(t, rs.Columns())
	require.Nil(t, rs.Row())
}
