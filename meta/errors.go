// --- #

package meta

import "github.com/pingcap/errors"

var (
	// ErrGraphExists is the error for db exists.
	ErrGraphExists = errors.New("graph exists")
	// ErrGraphNotExists is the error for db not exists.
	ErrGraphNotExists    = errors.New("graph not exists")
	ErrInvalidString     = errors.New("invalid string")
	ErrLabelExists       = errors.New("label exists")
	ErrLabelNotExists    = errors.New("label not exists")
	ErrIndexExists       = errors.New("index exists")
	ErrIndexNotExists    = errors.New("index not exists")
	ErrPropertyExists    = errors.New("property exists")
	ErrPropertyNotExists = errors.New("property not exists")
	ErrNoGraphSelected   = errors.New("no graph selected")
)
