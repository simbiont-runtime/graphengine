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

func TestProperty(t *testing.T) {
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

	propertyNames := []string{
		"property1",
		"property2",
		"property3",
		"property4",
		"property5",
	}

	// Create propertys
	for i, name := range propertyNames {
		err := kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
			meta := New(txn)
			info := &model.PropertyInfo{
				ID:   uint16(i) + 1,
				Name: model.NewCIStr(name),
			}
			return meta.CreateProperty(graphID, info)
		})
		assert.Nil(err)
	}

	// List properties
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		propertys, err := meta.ListProperties(graphID)
		assert.Nil(err)
		names := make([]string, 0, len(propertys))
		ids := make([]uint16, 0, len(propertys))
		for _, g := range propertys {
			names = append(names, g.Name.L)
			ids = append(ids, g.ID)
		}

		assert.Equal(propertyNames, names)
		assert.Equal([]uint16{1, 2, 3, 4, 5}, ids)

		return nil
	})
	assert.Nil(err)

	// Update property
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		err := meta.UpdateProperty(graphID, &model.PropertyInfo{
			ID:   4,
			Name: model.NewCIStr("PROPERTY4-modified"),
		})
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	// Get property
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		property, err := meta.GetProperty(graphID, 4)
		assert.Nil(err)
		assert.Equal(uint16(4), property.ID)
		assert.Equal("property4-modified", property.Name.L)
		return nil
	})
	assert.Nil(err)

	// Drop property
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		err := meta.DropProperty(graphID, 3)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	// List propertys again
	err = kv.TxnContext(context.TODO(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		propertys, err := meta.ListProperties(graphID)
		assert.Nil(err)
		names := make([]string, 0, len(propertys))
		ids := make([]uint16, 0, len(propertys))
		for _, g := range propertys {
			names = append(names, g.Name.L)
			ids = append(ids, g.ID)
		}

		propertyNames := []string{
			"property1",
			"property2",
			"property4-modified",
			"property5",
		}
		assert.Equal(propertyNames, names)
		assert.Equal([]uint16{1, 2, 4, 5}, ids)

		return nil
	})
	assert.Nil(err)
}
