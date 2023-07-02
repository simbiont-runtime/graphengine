// ---

package resolver

import "fmt"

// ErrRetryable suggests that client may restart the txn.
type ErrRetryable string

func (e ErrRetryable) Error() string {
	return fmt.Sprintf("retryable: %s", string(e))
}

// ErrAlreadyCommitted is returned specially when client tries to rollback a
// committed lock.
type ErrAlreadyCommitted uint64

func (e ErrAlreadyCommitted) Error() string {
	return "txn already committed"
}
