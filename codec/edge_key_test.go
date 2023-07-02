// ---

package codec

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEdgeKey(t *testing.T) {
	cases := []struct {
		graphID     int64
		srcVertexID int64
		dstVertexID int64
	}{
		{
			graphID:     100,
			srcVertexID: 200,
			dstVertexID: 300,
		},
		{
			graphID:     math.MaxInt64,
			srcVertexID: math.MaxInt64,
			dstVertexID: math.MaxInt64,
		},
	}
	for _, c := range cases {
		incomingEdgeKey := IncomingEdgeKey(c.graphID, c.srcVertexID, c.dstVertexID)
		graphID, srcVertexID, dstVertexID, err := ParseIncomingEdgeKey(incomingEdgeKey)
		require.NoError(t, err)
		require.Equal(t, c.graphID, graphID)
		require.Equal(t, c.srcVertexID, srcVertexID)
		require.Equal(t, c.dstVertexID, dstVertexID)

		outgoingEdgeKey := OutgoingEdgeKey(c.graphID, c.srcVertexID, c.dstVertexID)
		graphID, srcVertexID, dstVertexID, err = ParseOutgoingEdgeKey(outgoingEdgeKey)
		require.NoError(t, err)
		require.Equal(t, c.graphID, graphID)
		require.Equal(t, c.srcVertexID, srcVertexID)
		require.Equal(t, c.dstVertexID, dstVertexID)

		require.NotEqual(t, incomingEdgeKey, outgoingEdgeKey)
	}
}
