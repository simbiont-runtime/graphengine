// ---

package codec

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVertexKey(t *testing.T) {
	cases := []struct {
		graphID  int64
		vertexID int64
	}{
		{
			graphID:  100,
			vertexID: 200,
		}, {
			graphID:  math.MaxInt64,
			vertexID: math.MaxInt64,
		},
	}
	for _, c := range cases {
		key := VertexKey(c.graphID, c.vertexID)
		graphID, vertexID, err := ParseVertexKey(key)
		require.NoError(t, err)
		require.Equal(t, c.graphID, graphID)
		require.Equal(t, c.vertexID, vertexID)
	}
}
