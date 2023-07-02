// ---

package mvcc

type Op int32

const (
	Op_Put      Op = 0
	Op_Del      Op = 1
	Op_Lock     Op = 2
	Op_Rollback Op = 3
	// insert operation has a constraint that key should not exist before.
	Op_Insert         Op = 4
	Op_CheckNotExists Op = 5
)
