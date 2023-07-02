// ---

package mvcc

import (
	"fmt"

	"github.com/simbiont-runtime/graphengine/storage/kv"
)

// LockedError is returned when trying to Read/Write on a locked key. Caller should
// backoff or cleanup the lock then retry.
type LockedError struct {
	Key      kv.Key
	Primary  kv.Key
	StartVer kv.Version
	TTL      uint64
}

// Error formats the lock to a string.
func (e *LockedError) Error() string {
	return fmt.Sprintf("key is locked, key: %q, primary: %q, startVer: %v",
		e.Key, e.Primary, e.StartVer)
}
