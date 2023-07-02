// ---

package catalog

import (
	"testing"

	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/stretchr/testify/assert"
)

func TestCatalog_Apply(t *testing.T) {
	assert := assert.New(t)

	catalog := &Catalog{
		byID:   map[int64]*Graph{},
		byName: map[string]*Graph{},
	}

	cases := []struct {
		patch   *Patch
		checker func()
	}{
		{
			patch: &Patch{
				Type: PatchTypeCreateGraph,
				Data: &model.GraphInfo{
					ID:   1,
					Name: model.NewCIStr("graph1"),
				},
			},
			checker: func() {
				assert.NotNil(catalog.Graph("graph1"))
			},
		},
		{
			patch: &Patch{
				Type: PatchTypeCreateLabel,
				Data: &PatchLabel{
					GraphID: 1,
					LabelInfo: &model.LabelInfo{
						ID:   2,
						Name: model.NewCIStr("label1"),
					},
				},
			},
			checker: func() {
				graph := catalog.Graph("graph1")
				assert.NotNil(graph.Label("label1"))
			},
		},
		{
			patch: &Patch{
				Type: PatchTypeCreateLabel,
				Data: &PatchLabel{
					GraphID: 1,
					LabelInfo: &model.LabelInfo{
						ID:   3,
						Name: model.NewCIStr("label2"),
					},
				},
			},
			checker: func() {
				graph := catalog.Graph("graph1")
				assert.NotNil(graph.Label("label2"))
			},
		},
		{
			patch: &Patch{
				Type: PatchTypeDropLabel,
				Data: &PatchLabel{
					GraphID: 1,
					LabelInfo: &model.LabelInfo{
						ID:   2,
						Name: model.NewCIStr("label1"),
					},
				},
			},
			checker: func() {
				graph := catalog.Graph("graph1")
				assert.Nil(graph.Label("label1"))
				assert.NotNil(graph.Label("label2"))
			},
		},
		{
			patch: &Patch{
				Type: PatchTypeCreateProperties,
				Data: &PatchProperties{
					GraphID:   1,
					MaxPropID: 2,
					Properties: []*model.PropertyInfo{
						{
							ID:   1,
							Name: model.NewCIStr("property1"),
						},
						{
							ID:   2,
							Name: model.NewCIStr("property2"),
						},
					},
				},
			},
			checker: func() {
				graph := catalog.Graph("graph1")
				assert.NotNil(graph.Property("property1"))
				assert.NotNil(graph.Property("property2"))
			},
		},
		{
			patch: &Patch{
				Type: PatchTypeDropGraph,
				Data: &model.GraphInfo{
					ID:   1,
					Name: model.NewCIStr("graph1"),
				},
			},
			checker: func() {
				assert.Nil(catalog.Graph("graph1"))
			},
		},
	}

	for _, c := range cases {
		catalog.Apply(c.patch)
		if c.checker != nil {
			c.checker()
		}
	}
}
