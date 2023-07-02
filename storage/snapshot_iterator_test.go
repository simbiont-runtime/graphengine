// ---

package storage

import (
	"reflect"
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/simbiont-runtime/graphengine/storage/mvcc"
	"github.com/stretchr/testify/assert"
)

func TestIterator_Basic(t *testing.T) {
	s, err := Open(t.TempDir())
	assert.Nil(t, err)
	assert.NotNil(t, s)

	db := s.(*mvccStorage).db
	writes := db.NewBatch()

	// Mock data
	data := []struct {
		key []byte
		val []byte
		ver kv.Version
	}{
		{key: []byte("test"), val: []byte("test1"), ver: 100},
		{key: []byte("test"), val: []byte("test3"), ver: 300},
		{key: []byte("test1"), val: []byte("test5"), ver: 500},
		{key: []byte("test1"), val: []byte("test7"), ver: 700},
		{key: []byte("test2"), val: []byte("test9"), ver: 900},
	}
	wo := &pebble.WriteOptions{}
	for _, d := range data {
		v := mvcc.Value{
			Type:      mvcc.ValueTypePut,
			StartVer:  d.ver,
			CommitVer: d.ver + 10,
			Value:     d.val,
		}
		val, err := v.MarshalBinary()
		assert.Nil(t, err)
		assert.NotNil(t, val)
		err = writes.Set(mvcc.Encode(d.key, d.ver+1), val, wo)
		assert.Nil(t, err)
	}

	err = db.Apply(writes, wo)
	assert.Nil(t, err)

	type pair struct {
		k string
		v string
	}
	expected := []struct {
		ver     uint64
		results []pair
	}{
		{
			ver: 200,
			results: []pair{
				{"test", "test1"},
			},
		},
		{
			ver: 320,
			results: []pair{
				{"test", "test3"},
			},
		},
		{
			ver: 1000,
			results: []pair{
				{"test", "test3"},
				{"test1", "test7"},
				{"test2", "test9"},
			},
		},
	}
	for _, e := range expected {
		snapshot, err := s.Snapshot(kv.Version(e.ver))
		assert.Nil(t, err)
		iter, err := snapshot.Iter([]byte("t"), []byte("u"))
		assert.Nil(t, err)
		var results []pair

		for iter.Valid() {
			key := iter.Key()
			val := iter.Value()
			results = append(results, pair{string(key), string(val)})
			err = iter.Next()
			assert.Nil(t, err)
		}
		assert.True(t, reflect.DeepEqual(results, e.results), "expected: %v, got: %v (ver:%d)",
			e.results, results, e.ver)

		iter.Close()
	}
}

func TestReverseIterator(t *testing.T) {
	s, err := Open(t.TempDir())
	assert.Nil(t, err)
	assert.NotNil(t, s)

	db := s.(*mvccStorage).db
	writes := db.NewBatch()

	// Mock data
	data := []struct {
		key []byte
		val []byte
		ver kv.Version
	}{
		{key: []byte("test"), val: []byte("test1"), ver: 100},
		{key: []byte("test"), val: []byte("test3"), ver: 300},
		{key: []byte("test1"), val: []byte("test5"), ver: 500},
		{key: []byte("test1"), val: []byte("test7"), ver: 700},
		{key: []byte("test2"), val: []byte("test9"), ver: 900},
	}
	wo := &pebble.WriteOptions{}
	for _, d := range data {
		v := mvcc.Value{
			Type:      mvcc.ValueTypePut,
			StartVer:  d.ver,
			CommitVer: d.ver + 10,
			Value:     d.val,
		}
		val, err := v.MarshalBinary()
		assert.Nil(t, err)
		assert.NotNil(t, val)
		err = writes.Set(mvcc.Encode(d.key, d.ver+1), val, wo)
		assert.Nil(t, err)
	}

	err = db.Apply(writes, wo)
	assert.Nil(t, err)

	type pair struct {
		k string
		v string
	}
	expected := []struct {
		ver     uint64
		results []pair
	}{
		{
			ver: 200,
			results: []pair{
				{"test", "test1"},
			},
		},
		{
			ver: 320,
			results: []pair{
				{"test", "test3"},
			},
		},
		{
			ver: 1000,
			results: []pair{
				{"test2", "test9"},
				{"test1", "test7"},
				{"test", "test3"},
			},
		},
	}
	for _, e := range expected {
		snapshot, err := s.Snapshot(kv.Version(e.ver))
		assert.Nil(t, err)
		iter, err := snapshot.IterReverse([]byte("t"), []byte("u"))
		assert.Nil(t, err)
		var results []pair

		for iter.Valid() {
			key := iter.Key()
			val := iter.Value()
			results = append(results, pair{string(key), string(val)})
			err = iter.Next()
			assert.Nil(t, err)
		}
		assert.True(t, reflect.DeepEqual(results, e.results), "expected: %v, got: %v (ver:%d)",
			e.results, results, e.ver)

		iter.Close()
	}
}
