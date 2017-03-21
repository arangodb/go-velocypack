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

func TestSliceInt1(t *testing.T) {
	slice := velocypack.Slice{0x20, 0x33}
	value := int64(0x33)

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(2), slice.MustByteSize(), t)

	ASSERT_EQ(value, slice.MustGetInt(), t)
	ASSERT_EQ(value, slice.MustGetSmallInt(), t)
}

func TestSliceInt2(t *testing.T) {
	slice := velocypack.Slice{0x21, 0x23, 0x42}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(3), slice.MustByteSize(), t)

	ASSERT_EQ(int64(0x4223), slice.MustGetInt(), t)
	ASSERT_EQ(int64(0x4223), slice.MustGetSmallInt(), t)
}

func TestSliceInt3(t *testing.T) {
	slice := velocypack.Slice{0x22, 0x23, 0x42, 0x66}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(4), slice.MustByteSize(), t)

	ASSERT_EQ(int64(0x664223), slice.MustGetInt(), t)
	ASSERT_EQ(int64(0x664223), slice.MustGetSmallInt(), t)
}

func TestSliceInt4(t *testing.T) {
	slice := velocypack.Slice{0x23, 0x23, 0x42, 0x66, 0x7c}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(5), slice.MustByteSize(), t)

	ASSERT_EQ(int64(0x7c664223), slice.MustGetInt(), t)
	ASSERT_EQ(int64(0x7c664223), slice.MustGetSmallInt(), t)
}

func TestSliceInt5(t *testing.T) {
	slice := velocypack.Slice{0x24, 0x23, 0x42, 0x66, 0xac, 0x6f}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(6), slice.MustByteSize(), t)

	ASSERT_EQ(int64(0x6fac664223), slice.MustGetInt(), t)
	ASSERT_EQ(int64(0x6fac664223), slice.MustGetSmallInt(), t)
}

func TestSliceInt6(t *testing.T) {
	slice := velocypack.Slice{0x25, 0x23, 0x42, 0x66, 0xac, 0xff, 0x3f}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(7), slice.MustByteSize(), t)

	ASSERT_EQ(int64(0x3fffac664223), slice.MustGetInt(), t)
	ASSERT_EQ(int64(0x3fffac664223), slice.MustGetSmallInt(), t)
}

func TestSliceInt7(t *testing.T) {
	slice := velocypack.Slice{0x26, 0x23, 0x42, 0x66, 0xac, 0xff, 0x3f, 0x5a}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(8), slice.MustByteSize(), t)

	ASSERT_EQ(int64(0x5a3fffac664223), slice.MustGetInt(), t)
	ASSERT_EQ(int64(0x5a3fffac664223), slice.MustGetSmallInt(), t)
}

func TestSliceInt8(t *testing.T) {
	slice := velocypack.Slice{0x27, 0x23, 0x42, 0x66, 0xac, 0xff, 0x3f, 0xfa, 0x6f}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(9), slice.MustByteSize(), t)

	ASSERT_EQ(int64(0x6ffa3fffac664223), slice.MustGetInt(), t)
	ASSERT_EQ(int64(0x6ffa3fffac664223), slice.MustGetSmallInt(), t)
}

func TestSliceIntMax(t *testing.T) {
	t.Skip("TODO")
	/*	  Builder b;
	  b.add(Value(INT64_MAX));

	  Slice slice(b.slice());

		ASSERT_EQ(velocypack.Int, slice.Type(), t)
		ASSERT_TRUE(slice.IsInt(), t)
		ASSERT_EQ(velocypack.ValueLength(9), slice.MustByteSize(), t)

		ASSERT_EQ(int64(math.MaxInt64), slice.MustGetInt(), t)
	*/
}

func TestSliceNegInt1(t *testing.T) {
	slice := velocypack.Slice{0x20, 0xa3}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(2), slice.MustByteSize(), t)

	ASSERT_EQ(staticCastInt64(0xffffffffffffffa3), slice.MustGetInt(), t)
}

func TestSliceNegInt2(t *testing.T) {
	slice := velocypack.Slice{0x21, 0x23, 0xe2}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(3), slice.MustByteSize(), t)

	ASSERT_EQ(staticCastInt64(0xffffffffffffe223), slice.MustGetInt(), t)
}

func TestSliceNegInt3(t *testing.T) {
	slice := velocypack.Slice{0x22, 0x23, 0x42, 0xd6}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(4), slice.MustByteSize(), t)

	ASSERT_EQ(staticCastInt64(0xffffffffffd64223), slice.MustGetInt(), t)
}

func TestSliceNegInt4(t *testing.T) {
	slice := velocypack.Slice{0x23, 0x23, 0x42, 0x66, 0xac}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(5), slice.MustByteSize(), t)

	ASSERT_EQ(staticCastInt64(0xffffffffac664223), slice.MustGetInt(), t)
}

func TestSliceNegInt5(t *testing.T) {
	slice := velocypack.Slice{0x24, 0x23, 0x42, 0x66, 0xac, 0xff}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(6), slice.MustByteSize(), t)

	ASSERT_EQ(staticCastInt64(0xffffffffac664223), slice.MustGetInt(), t)
}

func TestSliceNegInt6(t *testing.T) {
	slice := velocypack.Slice{0x25, 0x23, 0x42, 0x66, 0xac, 0xff, 0xef}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(7), slice.MustByteSize(), t)

	ASSERT_EQ(staticCastInt64(0xffffefffac664223), slice.MustGetInt(), t)
}

func TestSliceNegInt7(t *testing.T) {
	slice := velocypack.Slice{0x26, 0x23, 0x42, 0x66, 0xac, 0xff, 0xef, 0xfa}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(8), slice.MustByteSize(), t)

	ASSERT_EQ(staticCastInt64(0xfffaefffac664223), slice.MustGetInt(), t)
}

func TestSliceNegInt8(t *testing.T) {
	slice := velocypack.Slice{0x27, 0x23, 0x42, 0x66, 0xac, 0xff, 0xef, 0xfa, 0x8e}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(9), slice.MustByteSize(), t)

	ASSERT_EQ(staticCastInt64(0x8efaefffac664223), slice.MustGetInt(), t)
}
