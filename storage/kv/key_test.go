//  Copyright 2022  GraphEngine Authors. All rights reserved.
//
// Copyright 2021 PingCAP, Inc.
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

package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixNextKey(t *testing.T) {
	k1 := []byte{0xff}
	k2 := []byte{0xff, 0xff}
	k3 := []byte{0xff, 0xff, 0xff, 0xff}

	pk1 := PrefixNextKey(k1)
	pk2 := PrefixNextKey(k2)
	pk3 := PrefixNextKey(k3)
	assert.Equal(t, []byte(""), pk1)
	assert.Equal(t, []byte(""), pk2)
	assert.Equal(t, []byte(""), pk3)
}
