// --- #

package meta

import (
	"context"
	"testing"

	"github.com/simbiont-runtime/graphengine/storage"
	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/stretchr/testify/assert"
)

func TestMeta_GlobalID(t *testing.T) {
	assert := assert.New(t)
	store, err := storage.Open(t.TempDir())
	assert.Nil(err)

	err = kv.TxnContext(context.Background(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		id, err := meta.GlobalID()
		assert.Nil(err)
		assert.Zero(id)
		return nil
	})
	assert.Nil(err)

	err = kv.TxnContext(context.Background(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		id, err := meta.NextGlobalID()
		assert.Nil(err)
		assert.Equal(int64(1), id)
		return nil
	})
	assert.Nil(err)

	err = kv.TxnContext(context.Background(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		id, err := meta.AdvanceGlobalID(10)
		assert.Nil(err)
		assert.Equal(int64(1), id)

		id, err = meta.GlobalID()
		assert.Nil(err)
		assert.Equal(int64(11), id)
		return nil
	})
	assert.Nil(err)

	err = kv.TxnContext(context.Background(), store, func(_ context.Context, txn kv.Transaction) error {
		meta := New(txn)
		ids, err := meta.GenGlobalIDs(3)
		assert.Nil(err)
		assert.Equal([]int64{12, 13, 14}, ids)
		return nil
	})
	assert.Nil(err)
}
