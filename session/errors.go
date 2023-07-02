// ---

package session

import "github.com/pingcap/errors"

var (
	ErrMultipleStatementsNotSupported = errors.New("multiple statements not supported")
	ErrFieldCountNotMatch             = errors.New("field count not match")
)
