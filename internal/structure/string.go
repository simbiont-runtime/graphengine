//  Copyright 2022  GraphEngine Authors. All rights reserved.
//
// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package structure

import (
	"context"
	"strconv"

	"github.com/pingcap/errors"
	"github.com/simbiont-runtime/graphengine/storage/kv"
)

// Set sets the string value of the key.
func (t *TxStructure) Set(key []byte, value []byte) error {
	if t.readWriter == nil {
		return ErrWriteOnSnapshot
	}
	ek := t.encodeStringDataKey(key)
	return t.readWriter.Set(ek, value)
}

// Get gets the string value of a key.
func (t *TxStructure) Get(key []byte) ([]byte, error) {
	ek := t.encodeStringDataKey(key)
	value, err := t.reader.Get(context.TODO(), ek)
	if errors.Cause(err) == kv.ErrNotExist {
		err = nil
	}
	return value, errors.Trace(err)
}

// GetInt64 gets the int64 value of a key.
func (t *TxStructure) GetInt64(key []byte) (int64, error) {
	v, err := t.Get(key)
	if err != nil || v == nil {
		return 0, errors.Trace(err)
	}

	n, err := strconv.ParseInt(string(v), 10, 64)
	return n, errors.Trace(err)
}

// Inc increments the integer value of a key by step, returns
// the value after the increment.
func (t *TxStructure) Inc(key []byte, step int64) (int64, error) {
	if t.readWriter == nil {
		return 0, ErrWriteOnSnapshot
	}
	ek := t.encodeStringDataKey(key)
	// txn Inc will lock this key, so we don't lock it here.
	n, err := IncInt64(t.readWriter, ek, step)
	if errors.Cause(err) == kv.ErrNotExist {
		err = nil
	}
	return n, errors.Trace(err)
}

// Clear removes the string value of the key.
func (t *TxStructure) Clear(key []byte) error {
	if t.readWriter == nil {
		return ErrWriteOnSnapshot
	}
	ek := t.encodeStringDataKey(key)
	err := t.readWriter.Delete(ek)
	if errors.Cause(err) == kv.ErrNotExist {
		err = nil
	}
	return errors.Trace(err)
}

// IncInt64 increases the value for key k in kv store by step.
func IncInt64(rm kv.RetrieverMutator, k kv.Key, step int64) (int64, error) {
	val, err := rm.Get(context.TODO(), k)
	if errors.Cause(err) == kv.ErrNotExist {
		err = rm.Set(k, []byte(strconv.FormatInt(step, 10)))
		if err != nil {
			return 0, err
		}
		return step, nil
	}
	if err != nil {
		return 0, err
	}

	intVal, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		return 0, errors.Trace(err)
	}

	intVal += step
	err = rm.Set(k, []byte(strconv.FormatInt(intVal, 10)))
	if err != nil {
		return 0, err
	}
	return intVal, nil
}
