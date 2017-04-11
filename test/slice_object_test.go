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

func TestSliceObjectEmpty(t *testing.T) {
	slice := velocypack.Slice{0x0a}
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_TRUE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(1), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(0), mustLength(slice.Length()), t)
}

func TestSliceObjectCases1(t *testing.T) {
	slice := velocypack.Slice{0x0b, 0x00, 0x03, 0x41, 0x61, 0x31, 0x41, 0x62,
		0x32, 0x41, 0x63, 0x33, 0x03, 0x06, 0x09}
	slice[1] = byte(len(slice)) // Set byte length
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_FALSE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(len(slice)), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.Length()), t)
	ss := mustSlice(slice.Get("a"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(1), mustInt(ss.GetInt()), t)
	ASSERT_EQ(int64(1), mustInt(mustSlice(slice.ValueAt(0)).GetInt()), t)
}

func TestSliceObjectCases2(t *testing.T) {
	slice := velocypack.Slice{0x0b, 0x00, 0x03, 0x00, 0x00, 0x41, 0x61, 0x31, 0x41,
		0x62, 0x32, 0x41, 0x63, 0x33, 0x05, 0x08, 0x0b}
	slice[1] = byte(len(slice)) // Set byte length
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_FALSE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(len(slice)), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.Length()), t)
	ss := mustSlice(slice.Get("a"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(1), mustInt(ss.GetInt()), t)
}

func TestSliceObjectCases3(t *testing.T) {
	slice := velocypack.Slice{0x0b, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x41, 0x61, 0x31, 0x41, 0x62,
		0x32, 0x41, 0x63, 0x33, 0x09, 0x0c, 0x0f}
	slice[1] = byte(len(slice)) // Set byte length
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_FALSE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(len(slice)), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.Length()), t)
	ss := mustSlice(slice.Get("a"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(1), mustInt(ss.GetInt()), t)
}

func TestSliceObjectCases7(t *testing.T) {
	slice := velocypack.Slice{0x0c, 0x00, 0x00, 0x03, 0x00, 0x41, 0x61, 0x31, 0x41, 0x62,
		0x32, 0x41, 0x63, 0x33, 0x05, 0x00, 0x08, 0x00, 0x0b, 0x00}
	slice[1] = byte(len(slice)) // Set byte length
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_FALSE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(len(slice)), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.Length()), t)
	ss := mustSlice(slice.Get("a"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(1), mustInt(ss.GetInt()), t)
}

func TestSliceObjectCases8(t *testing.T) {
	slice := velocypack.Slice{0x0c, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x41, 0x61, 0x31, 0x41, 0x62, 0x32, 0x41,
		0x63, 0x33, 0x09, 0x00, 0x0c, 0x00, 0x0f, 0x00}
	slice[1] = byte(len(slice)) // Set byte length
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_FALSE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(len(slice)), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.Length()), t)
	ss := mustSlice(slice.Get("a"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(1), mustInt(ss.GetInt()), t)
}

func TestSliceObjectCases11(t *testing.T) {
	slice := velocypack.Slice{0x0d, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x41,
		0x61, 0x31, 0x41, 0x62, 0x32, 0x41, 0x63, 0x33, 0x09, 0x00,
		0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00}
	slice[1] = byte(len(slice)) // Set byte length
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_FALSE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(len(slice)), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.Length()), t)
	ss := mustSlice(slice.Get("a"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(1), mustInt(ss.GetInt()), t)
}

func TestSliceObjectCases13(t *testing.T) {
	slice := velocypack.Slice{0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x41,
		0x61, 0x31, 0x41, 0x62, 0x32, 0x41, 0x63, 0x33, 0x09, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	slice[1] = byte(len(slice)) // Set byte length
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_FALSE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(len(slice)), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.Length()), t)
	ss := mustSlice(slice.Get("a"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(1), mustInt(ss.GetInt()), t)
}

func TestSliceObjectCompact(t *testing.T) {
	slice := velocypack.Slice{0x14, 0x0f, 0x41, 0x61, 0x30, 0x41, 0x62, 0x31,
		0x41, 0x63, 0x32, 0x41, 0x64, 0x33, 0x04}
	slice[1] = byte(len(slice)) // Set byte length
	assertEqualFromReader(t, slice)

	ASSERT_EQ(velocypack.Object, slice.Type(), t)
	ASSERT_TRUE(slice.IsObject(), t)
	ASSERT_FALSE(slice.IsEmptyObject(), t)
	ASSERT_EQ(velocypack.ValueLength(len(slice)), mustLength(slice.ByteSize()), t)
	ASSERT_EQ(velocypack.ValueLength(4), mustLength(slice.Length()), t)
	ss := mustSlice(slice.Get("a"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(0), mustInt(ss.GetInt()), t)

	ss = mustSlice(slice.Get("b"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(1), mustInt(ss.GetInt()), t)

	ss = mustSlice(slice.Get("d"))
	ASSERT_TRUE(ss.IsSmallInt(), t)
	ASSERT_EQ(int64(3), mustInt(ss.GetInt()), t)
}
