// ---

package compiler

import "github.com/pingcap/errors"

var (
	ErrIncorrectGraphName        = errors.New("incorrect graph name")
	ErrIncorrectLabelName        = errors.New("incorrect label name")
	ErrIncorrectIndexName        = errors.New("incorrect index name")
	ErrGraphNotChosen            = errors.New("please choose graph first")
	ErrVariableReferenceNotExits = errors.New("reference not exists variable")
)
