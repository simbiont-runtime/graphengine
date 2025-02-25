//  Copyright 2022  GraphEngine Authors. All rights reserved.
//
// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package latch

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/simbiont-runtime/graphengine/storage/kv"
	"github.com/sourcegraph/conc"
)

func TestWithConcurrency(t *testing.T) {
	sched := NewScheduler(7)
	sched.Run()
	defer sched.Close()
	rand.Seed(time.Now().Unix())

	ch := make(chan []kv.Key, 100)
	const workerCount = 10
	var wg conc.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Go(func() {
			for txn := range ch {
				lock := sched.Lock(getTso(), txn)
				if lock.IsStale() {
					// Should restart the transaction or return error
				} else {
					lock.SetCommitVer(getTso())
					// Do 2pc
				}
				sched.UnLock(lock)
			}
		})
	}

	for i := 0; i < 999; i++ {
		ch <- generate()
	}
	close(ch)

	wg.Wait()
}

// generate generates something like:
// {[]byte("a"), []byte("b"), []byte("c")}
// {[]byte("a"), []byte("d"), []byte("e"), []byte("f")}
// {[]byte("e"), []byte("f"), []byte("g"), []byte("h")}
// The data should not repeat in the sequence.
func generate() []kv.Key {
	table := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	ret := make([]kv.Key, 0, 5)
	chance := []int{100, 60, 40, 20}
	for i := 0; i < len(chance); i++ {
		needMore := rand.Intn(100) < chance[i]
		if needMore {
			randBytes := []byte{table[rand.Intn(len(table))]}
			if !contains(randBytes, ret) {
				ret = append(ret, randBytes)
			}
		}
	}
	return ret
}

func contains(x []byte, set []kv.Key) bool {
	for _, y := range set {
		if bytes.Equal(x, y) {
			return true
		}
	}
	return false
}
