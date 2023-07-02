// ---

package codec

import (
	"testing"

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/stretchr/testify/assert"
)

func TestPropertyEncoder_Encode(t *testing.T) {
	cases := []struct {
		labelIDs    []uint16
		propertyIDs []uint16
		values      []datum.Datum
	}{
		{
			labelIDs:    []uint16{1, 2, 3},
			propertyIDs: []uint16{1, 2, 3},
			values: []datum.Datum{
				datum.NewString("hello"),
				datum.NewInt(1),
				datum.NewFloat(1.1),
			},
		},
	}

	for _, c := range cases {
		encoder := &PropertyEncoder{}
		// FIXME: validate the values.
		_, err := encoder.Encode(nil, c.labelIDs, c.propertyIDs, c.values)
		assert.NoError(t, err)
	}
}
