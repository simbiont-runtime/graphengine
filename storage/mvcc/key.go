// ---

package mvcc

import "github.com/simbiont-runtime/graphengine/codec"

// Key is the encoded key type with timestamp.
type Key []byte

// NewKey encodes a key into MvccKey.
func NewKey(key []byte) Key {
	if len(key) == 0 {
		return nil
	}
	return codec.EncodeBytes(nil, key)
}

// Raw decodes a MvccKey to original key.
func (key Key) Raw() []byte {
	if len(key) == 0 {
		return nil
	}
	_, k, err := codec.DecodeBytes(key, nil)
	if err != nil {
		panic(err)
	}
	return k
}
