// ---

package graphengine

const defaultConcurrency = 512

//	Options contains some options which is used to customize the  GraphEngine database
//
// instance while instantiating graphengine.
type Options struct {
	// Concurrency is used to limit the max concurrent sessions count. The NewSession
	// method will block if the current alive sessions count reach this limitation.
	Concurrency int64
}

// SetDefaults sets the missing options into default value.
func (opt *Options) SetDefaults() {
	if opt.Concurrency <= 0 {
		opt.Concurrency = defaultConcurrency
	}
}
