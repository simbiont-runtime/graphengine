// ---

package codec

import (
	"fmt"
	"testing"

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/stretchr/testify/assert"
)

func TestPropertyDecoder_Decode(t *testing.T) {
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
		{
			labelIDs:    []uint16{2, 3, 1},
			propertyIDs: []uint16{2, 3, 1},
			values: []datum.Datum{
				datum.NewInt(1),
				datum.NewFloat(1.1),
				datum.NewString("hello"),
			},
		},
	}

	for _, c := range cases {
		encoder := &PropertyEncoder{}
		bytes, err := encoder.Encode(nil, c.labelIDs, c.propertyIDs, c.values)
		assert.Nil(t, err)

		var labels []*model.LabelInfo
		for _, id := range c.labelIDs {
			labels = append(labels, &model.LabelInfo{
				ID:   int64(id),
				Name: model.NewCIStr(fmt.Sprintf("label%d", id)),
			})
		}

		var properties []*model.PropertyInfo
		for _, id := range c.propertyIDs {
			properties = append(properties, &model.PropertyInfo{
				ID:   id,
				Name: model.NewCIStr(fmt.Sprintf("property%d", id)),
			})
		}

		decoder := NewPropertyDecoder(labels, properties)
		labelIDs, row, err := decoder.Decode(bytes)
		assert.NoError(t, err)
		assert.Equal(t, map[uint16]struct{}{
			1: {},
			2: {},
			3: {},
		}, labelIDs)
		assert.Equal(t, map[uint16]datum.Datum{
			1: datum.NewString("hello"),
			2: datum.NewInt(1),
			3: datum.NewFloat(1.1),
		}, row)
	}
}
