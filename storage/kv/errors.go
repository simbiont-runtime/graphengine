// ---

package kv

import (
	"fmt"

	"github.com/pingcap/errors"
)

var (
	// ErrTxnConflicts indicates the current transaction contains some vertex/edge/index
	// conflicts with others.
	ErrTxnConflicts = errors.New("transaction conflicts")

	// ErrNotExist means the related data not exist.
	ErrNotExist = errors.New("not exist")

	// ErrCannotSetNilValue is the error when sets an empty value.
	ErrCannotSetNilValue = errors.New("can not set nil value")

	// ErrInvalidTxn is the error when commits or rollbacks in an invalid transaction.
	ErrInvalidTxn = errors.New("invalid transaction")

	ErrInvalidStartVer = errors.New("invalid start timestamp for transaction")
)

// ErrEntryTooLarge is the error when a key value entry is too large.
type ErrEntryTooLarge struct {
	Limit uint64
	Size  uint64
}

func (e *ErrEntryTooLarge) Error() string {
	return fmt.Sprintf("entry size too large, size: %v,limit: %v.", e.Size, e.Limit)
}

// ErrTxnTooLarge is the error when transaction is too large, lock time reached the maximum value.
type ErrTxnTooLarge struct {
	Size int
}

func (e *ErrTxnTooLarge) Error() string {
	return fmt.Sprintf("txn too large, size: %v.", e.Size)
}

func IsErrNotFound(err error) bool {
	return errors.Cause(err) == ErrNotExist
}

// ErrKeyAlreadyExist is returned when key exists but this key has a constraint that
// it should not exist. Client should return duplicated entry error.
type ErrKeyAlreadyExist struct {
	Key []byte
}

func (e *ErrKeyAlreadyExist) Error() string {
	return fmt.Sprintf("key already exist, key: %q", e.Key)
}

// ErrGroup is used to collect multiple errors.
type ErrGroup struct {
	Errors []error
}

func (e *ErrGroup) Error() string {
	return fmt.Sprintf("Errors: %v", e.Errors)
}

// ErrConflict is returned when the commitTS of key in the DB is greater than startTS.
type ErrConflict struct {
	StartVer          Version
	ConflictStartVer  Version
	ConflictCommitVer Version
	Key               Key
}

func (e *ErrConflict) Error() string {
	return "write conflict"
}
