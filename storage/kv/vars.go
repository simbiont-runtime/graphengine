//  Copyright 2022  GraphEngine Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kv

import "sync/atomic"

// DefTxnCommitBatchSize is the default value of TxnCommitBatchSize.
const DefTxnCommitBatchSize uint64 = 16 * 1024

// TxnCommitBatchSize controls the batch size of transaction commit related requests sent by client to TiKV,
// TiKV recommends each RPC packet should be less than ~1MB.
var TxnCommitBatchSize atomic.Uint64

func init() {
	TxnCommitBatchSize.Store(DefTxnCommitBatchSize)
}
