// ---

package gc

import "github.com/simbiont-runtime/graphengine/storage/kv"

// Task represents a GC task which contains the key ranges to execute GC.
type Task struct {
	LowerBound kv.Key
	UpperBound kv.Key
}
