// ---

package resolver

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/simbiont-runtime/graphengine/storage/mvcc"
	"github.com/stretchr/testify/assert"
)

func TestNewScheduler(t *testing.T) {
	assert := assert.New(t)

	db, err := pebble.Open(t.TempDir(), nil)
	assert.Nil(err)

	s := NewScheduler(2)
	assert.NotNil(s)

	s.SetDB(db)
	s.Run()

	// Put some locks into database
	data := []struct {
		key        string
		primaryKey string
		startVer   kv.Version
		commitVer  kv.Version
	}{
		// Transaction A
		{
			key:        "a",
			primaryKey: "a",
			startVer:   1,
			commitVer:  2,
		},
		{
			key:        "b",
			primaryKey: "a",
			startVer:   1,
			commitVer:  2,
		},
		{
			key:        "c",
			primaryKey: "a",
			startVer:   1,
			commitVer:  2,
		},
	}

	batch := db.NewBatch()
	var keys []kv.Key
	for _, d := range data {
		lock := mvcc.Lock{
			StartVer: d.startVer,
			Primary:  kv.Key(d.primaryKey),
			Value:    []byte(d.key + "_value"),
			Op:       mvcc.Op_Put,
			TTL:      100,
		}

		writeVal, err := lock.MarshalBinary()
		assert.Nil(err)

		err = batch.Set(mvcc.Encode(kv.Key(d.key), mvcc.LockVer), writeVal, nil)
		assert.Nil(err)

		keys = append(keys, kv.Key(d.key))
	}

	err = batch.Commit(nil)
	assert.Nil(err)

	notifier := NewMultiKeysNotifier(len(data))
	s.Resolve(keys, 1, 2, notifier)
	errs := notifier.Wait()

	assert.Nil(errs)

	// Check the values
	for _, d := range data {
		val, _, err := db.Get(mvcc.Encode(kv.Key(d.key), d.commitVer))
		assert.Nil(err)
		v := mvcc.Value{}
		err = v.UnmarshalBinary(val)
		assert.Nil(err)
		assert.Equal(d.key+"_value", string(v.Value))
	}
}

func TestNewScheduler_EarlyResolve(t *testing.T) {
	assert := assert.New(t)

	db, err := pebble.Open(t.TempDir(), nil)
	assert.Nil(err)

	s := NewScheduler(2)
	assert.NotNil(s)

	s.SetDB(db)
	s.Run()

	// Put some locks into database
	data := []struct {
		key        string
		primaryKey string
		startVer   kv.Version
		commitVer  kv.Version
	}{
		// Transaction A
		{
			key:        "a",
			primaryKey: "a",
			startVer:   1,
			commitVer:  2,
		},
		{
			key:        "b",
			primaryKey: "a",
			startVer:   1,
			commitVer:  2,
		},
		{
			key:        "c",
			primaryKey: "a",
			startVer:   1,
			commitVer:  2,
		},
	}

	batch := db.NewBatch()
	var keys []kv.Key
	for _, d := range data {
		lock := mvcc.Lock{
			StartVer: d.startVer,
			Primary:  kv.Key(d.primaryKey),
			Value:    []byte(d.key + "_value"),
			Op:       mvcc.Op_Put,
			TTL:      100,
		}

		writeVal, err := lock.MarshalBinary()
		assert.Nil(err)

		err = batch.Set(mvcc.LockKey(kv.Key(d.key)), writeVal, nil)
		assert.Nil(err)

		keys = append(keys, kv.Key(d.key))
	}

	err = batch.Commit(nil)
	assert.Nil(err)

	// Resolve one of them
	resolveBatch := db.NewBatch()
	err = Resolve(db, resolveBatch, kv.Key("b"), 1, 2)
	assert.Nil(err)
	err = resolveBatch.Commit(nil)
	assert.Nil(err)
	val, _, err := db.Get(mvcc.Encode(kv.Key("b"), 2))
	v := mvcc.Value{}
	err = v.UnmarshalBinary(val)
	assert.Nil(err)
	assert.Equal("b_value", string(v.Value))

	// Try to resolve all keys and part of them were resolved.
	notifier := NewMultiKeysNotifier(len(data))
	s.Resolve(keys, 1, 2, notifier)
	errs := notifier.Wait()

	assert.Nil(errs)

	// Check the values
	for _, d := range data {
		val, _, err := db.Get(mvcc.Encode(kv.Key(d.key), d.commitVer))
		assert.Nil(err, "key:%s", d.key)
		v := mvcc.Value{}
		err = v.UnmarshalBinary(val)
		assert.Nil(err)
		assert.Equal(d.key+"_value", string(v.Value))
	}
}
