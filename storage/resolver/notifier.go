// ---

package resolver

import "sync"

// Notifier is used to notify the task finished.
type Notifier interface {
	Notify(err error)
	Wait() []error
}

type MultiKeysNotifier struct {
	wg   sync.WaitGroup
	mu   sync.RWMutex
	errs []error
}

func NewMultiKeysNotifier(size int) Notifier {
	n := &MultiKeysNotifier{}
	n.wg.Add(size)
	return n
}

func (s *MultiKeysNotifier) Notify(err error) {
	s.wg.Done()
	if err != nil {
		s.errs = append(s.errs, err)
	}
}

func (s *MultiKeysNotifier) Wait() []error {
	s.wg.Wait()
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errs
}
