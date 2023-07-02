// --- #

package meta

import (
	"context"
	"testing"

	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/simbiont-runtime/graphengine/storage"
	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	assert := assert.New(t)
	store, err := storage.Open(t.TempDir())
	assert.Nil(err)

	graphNames := []string{
		"graph1",
		"graph2",
		"graph3",
		"graph4",
		"graph5",
	}

	// Create graphs
	for i, name := range graphNames {
		err := kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
			meta := New(txn)
			info := &model.GraphInfo{
				ID:   int64(i) + 1,
				Name: model.NewCIStr(name),
			}
			return meta.CreateGraph(info)
		})
		assert.Nil(err)
	}

	// List graphs
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		graphs, err := meta.ListGraphs()
		assert.Nil(err)
		names := make([]string, 0, len(graphs))
		ids := make([]int64, 0, len(graphs))
		for _, g := range graphs {
			names = append(names, g.Name.L)
			ids = append(ids, g.ID)
		}

		assert.Equal(graphNames, names)
		assert.Equal([]int64{1, 2, 3, 4, 5}, ids)

		return nil
	})
	assert.Nil(err)

	// Update graph
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		err := meta.UpdateGraph(&model.GraphInfo{
			ID:   4,
			Name: model.NewCIStr("GRAPH4-modified"),
		})
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	// Get graph
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		graph, err := meta.GetGraph(4)
		assert.Nil(err)
		assert.Equal(int64(4), graph.ID)
		assert.Equal("graph4-modified", graph.Name.L)
		return nil
	})
	assert.Nil(err)

	// Drop graph
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		err := meta.DropGraph(3)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	// List graphs again
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		graphs, err := meta.ListGraphs()
		assert.Nil(err)
		names := make([]string, 0, len(graphs))
		ids := make([]int64, 0, len(graphs))
		for _, g := range graphs {
			names = append(names, g.Name.L)
			ids = append(ids, g.ID)
		}

		graphNames := []string{
			"graph1",
			"graph2",
			"graph4-modified",
			"graph5",
		}
		assert.Equal(graphNames, names)
		assert.Equal([]int64{1, 2, 4, 5}, ids)

		return nil
	})
	assert.Nil(err)
}
