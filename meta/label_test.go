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

func TestLabel(t *testing.T) {
	assert := assert.New(t)
	store, err := storage.Open(t.TempDir())
	assert.Nil(err)

	var graphID = int64(100)

	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		info := &model.GraphInfo{
			ID:   graphID,
			Name: model.NewCIStr("test-graph"),
		}
		return meta.CreateGraph(info)
	})
	assert.Nil(err)

	labelNames := []string{
		"label1",
		"label2",
		"label3",
		"label4",
		"label5",
	}

	// Create labels
	for i, name := range labelNames {
		err := kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
			meta := New(txn)
			info := &model.LabelInfo{
				ID:   int64(i) + 1,
				Name: model.NewCIStr(name),
			}
			return meta.CreateLabel(graphID, info)
		})
		assert.Nil(err)
	}

	// List labels
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		labels, err := meta.ListLabels(graphID)
		assert.Nil(err)
		names := make([]string, 0, len(labels))
		ids := make([]int64, 0, len(labels))
		for _, g := range labels {
			names = append(names, g.Name.L)
			ids = append(ids, g.ID)
		}

		assert.Equal(labelNames, names)
		assert.Equal([]int64{1, 2, 3, 4, 5}, ids)

		return nil
	})
	assert.Nil(err)

	// Update label
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		err := meta.UpdateLabel(graphID, &model.LabelInfo{
			ID:   4,
			Name: model.NewCIStr("LABEL4-modified"),
		})
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	// Get label
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		label, err := meta.GetLabel(graphID, 4)
		assert.Nil(err)
		assert.Equal(int64(4), label.ID)
		assert.Equal("label4-modified", label.Name.L)
		return nil
	})
	assert.Nil(err)

	// Drop label
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		err := meta.DropLabel(graphID, 3)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	// List labels again
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		labels, err := meta.ListLabels(graphID)
		assert.Nil(err)
		names := make([]string, 0, len(labels))
		ids := make([]int64, 0, len(labels))
		for _, g := range labels {
			names = append(names, g.Name.L)
			ids = append(ids, g.ID)
		}

		labelNames := []string{
			"label1",
			"label2",
			"label4-modified",
			"label5",
		}
		assert.Equal(labelNames, names)
		assert.Equal([]int64{1, 2, 4, 5}, ids)

		return nil
	})
	assert.Nil(err)
}
