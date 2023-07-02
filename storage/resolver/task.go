// ---

package resolver

import (
	"github.com/simbiont-runtime/graphengine/storage/kv"
)

// Task represents a resolve task.
type Task struct {
	Key       kv.Key
	StartVer  kv.Version
	CommitVer kv.Version
	Notifier  Notifier
}
