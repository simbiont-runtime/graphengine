// ---

package resolver

import (
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/pingcap/errors"
	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/simbiont-runtime/graphengine/storage/mvcc"
)

type TxnAction byte

const (
	TxnActionNone TxnAction = iota
	TxnActionTTLExpireRollback
	TxnActionLockNotExistRollback
	TxnActionLockNotExistDoNothing
)

type TxnStatus struct {
	CommitVer kv.Version
	Action    TxnAction
}

// CheckTxnStatus checks the transaction status according to the primary key.
func CheckTxnStatus(db *pebble.DB, vp kv.VersionProvider, primaryKey kv.Key, startVer kv.Version) (TxnStatus, error) {
	opt := pebble.IterOptions{LowerBound: mvcc.LockKey(primaryKey)}
	iter := db.NewIter(&opt)
	iter.First()
	defer iter.Close()

	if !iter.Valid() {
		return TxnStatus{}, errors.New("txn not found")
	}

	decoder := mvcc.LockDecoder{ExpectKey: primaryKey}
	exists, err := decoder.Decode(iter)
	if err != nil {
		return TxnStatus{}, err
	}

	// If the transaction lock exists means the current transaction not committed.
	if exists && decoder.Lock.StartVer == startVer {
		ver := vp.CurrentVersion()
		exp := startVer + kv.Version(time.Duration(decoder.Lock.TTL)*time.Millisecond)
		if exp < ver {
			return TxnStatus{Action: TxnActionTTLExpireRollback}, nil
		}
		return TxnStatus{Action: TxnActionNone}, nil
	}

	c, exists, err := getTxnCommitInfo(iter, primaryKey, startVer)
	if err != nil {
		return TxnStatus{}, err
	}
	if exists {
		if c.Type == mvcc.ValueTypeRollback {
			return TxnStatus{Action: TxnActionLockNotExistRollback}, nil
		}
		return TxnStatus{CommitVer: c.CommitVer, Action: TxnActionLockNotExistDoNothing}, nil
	}

	return TxnStatus{}, errors.New("transaction status missing")
}
