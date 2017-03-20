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

func TestSliceArrayEmpty(t *testing.T) {
	slice := velocypack.Slice{0x01}

	ASSERT_EQ(velocypack.Array, slice.Type(), t)
	ASSERT_TRUE(slice.IsArray(), t)
	ASSERT_TRUE(slice.IsEmptyArray(), t)
	ASSERT_EQ(velocypack.ValueLength(1), slice.MustByteSize(), t)
	ASSERT_EQ(velocypack.ValueLength(0), slice.MustLength(), t)
}
