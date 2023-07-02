// ---

package catalog

import (
	"sync/atomic"
	"testing"

	"github.com/simbiont-runtime/graphengine/meta"
	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/simbiont-runtime/graphengine/storage"
	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/stretchr/testify/assert"
)

func Test_Load(t *testing.T) {
	assert := assert.New(t)
	store, err := storage.Open(t.TempDir())
	assert.Nil(err)
	defer store.Close()

	ID := atomic.Int64{}
	cases := []*model.GraphInfo{
		{
			ID:   ID.Add(1),
			Name: model.NewCIStr("graph1"),
			Labels: []*model.LabelInfo{
				{
					ID:   ID.Add(1),
					Name: model.NewCIStr("label1"),
				},
				{
					ID:   ID.Add(1),
					Name: model.NewCIStr("label2"),
				},
			},
		},
		{
			ID:   ID.Add(1),
			Name: model.NewCIStr("graph2"),
			Labels: []*model.LabelInfo{
				{
					ID:   ID.Add(1),
					Name: model.NewCIStr("label1"),
				},
				{
					ID:   ID.Add(1),
					Name: model.NewCIStr("label2"),
				},
			},
			Properties: []*model.PropertyInfo{
				{
					ID:   uint16(ID.Add(1)),
					Name: model.NewCIStr("property1"),
				},
				{
					ID:   uint16(ID.Add(1)),
					Name: model.NewCIStr("property2"),
				},
			},
		},
	}

	// Create mock data.
	err = kv.Txn(store, func(txn kv.Transaction) error {
		meta := meta.New(txn)
		for _, g := range cases {
			err := meta.CreateGraph(g)
			assert.Nil(err)
			for _, l := range g.Labels {
				err := meta.CreateLabel(g.ID, l)
				assert.Nil(err)
			}
			for _, p := range g.Properties {
				err := meta.CreateProperty(g.ID, p)
				assert.Nil(err)
			}
		}
		return nil
	})
	assert.Nil(err)

	snapshot, err := store.Snapshot(store.CurrentVersion())
	assert.Nil(err)

	catalog, err := Load(snapshot)
	assert.Nil(err)
	assert.Equal(len(cases), len(catalog.byName))
	assert.Equal(len(catalog.byID), len(catalog.byName))

	for _, g := range cases {
		graph := catalog.Graph(g.Name.L)
		assert.Equal(g.ID, graph.Meta().ID)
		assert.Equal(graph, catalog.GraphByID(g.ID))

		// Labels
		for _, l := range g.Labels {
			label := graph.Label(l.Name.L)
			assert.Equal(l, label.Meta())
			assert.Equal(label, graph.LabelByID(l.ID))
		}

		// Properties
		for _, p := range g.Properties {
			property := graph.Property(p.Name.L)
			assert.Equal(p, property)
			assert.Equal(property, graph.PropertyByID(p.ID))
		}
	}
}
