//
// DISCLAIMER
//
// Copyright 2017 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package test

import (
	"testing"

	velocypack "github.com/arangodb/go-velocypack"
)

func TestSliceFalse(t *testing.T) {
	slice := velocypack.Slice{0x19}

	ASSERT_EQ(velocypack.Bool, slice.Type(), t)
	ASSERT_TRUE(slice.IsBool(), t)
	ASSERT_TRUE(slice.IsFalse(), t)
	ASSERT_FALSE(slice.IsTrue(), t)
	ASSERT_EQ(velocypack.ValueLength(1), slice.MustByteSize(), t)
	ASSERT_FALSE(slice.MustGetBool(), t)
}

func TestSliceTrue(t *testing.T) {
	slice := velocypack.Slice{0x1a}

	ASSERT_EQ(velocypack.Bool, slice.Type(), t)
	ASSERT_TRUE(slice.IsBool(), t)
	ASSERT_FALSE(slice.IsFalse(), t)
	ASSERT_TRUE(slice.IsTrue(), t)
	ASSERT_EQ(velocypack.ValueLength(1), slice.MustByteSize(), t)
	ASSERT_TRUE(slice.MustGetBool(), t)
}