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
	ASSERT_EQ(velocypack.ValueLength(2), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(value, mustInt(slice.GetInt()), t)
	ASSERT_EQ(value, mustInt(slice.GetSmallInt()), t)
}

func TestSliceInt2(t *testing.T) {
	slice := velocypack.Slice{0x21, 0x23, 0x42}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(int64(0x4223), mustInt(slice.GetInt()), t)
	ASSERT_EQ(int64(0x4223), mustInt(slice.GetSmallInt()), t)
}

func TestSliceInt3(t *testing.T) {
	slice := velocypack.Slice{0x22, 0x23, 0x42, 0x66}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(4), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(int64(0x664223), mustInt(slice.GetInt()), t)
	ASSERT_EQ(int64(0x664223), mustInt(slice.GetSmallInt()), t)
}

func TestSliceInt4(t *testing.T) {
	slice := velocypack.Slice{0x23, 0x23, 0x42, 0x66, 0x7c}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(5), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(int64(0x7c664223), mustInt(slice.GetInt()), t)
	ASSERT_EQ(int64(0x7c664223), mustInt(slice.GetSmallInt()), t)
}

func TestSliceInt5(t *testing.T) {
	slice := velocypack.Slice{0x24, 0x23, 0x42, 0x66, 0xac, 0x6f}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(6), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(int64(0x6fac664223), mustInt(slice.GetInt()), t)
	ASSERT_EQ(int64(0x6fac664223), mustInt(slice.GetSmallInt()), t)
}

func TestSliceInt6(t *testing.T) {
	slice := velocypack.Slice{0x25, 0x23, 0x42, 0x66, 0xac, 0xff, 0x3f}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(7), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(int64(0x3fffac664223), mustInt(slice.GetInt()), t)
	ASSERT_EQ(int64(0x3fffac664223), mustInt(slice.GetSmallInt()), t)
}

func TestSliceInt7(t *testing.T) {
	slice := velocypack.Slice{0x26, 0x23, 0x42, 0x66, 0xac, 0xff, 0x3f, 0x5a}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(8), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(int64(0x5a3fffac664223), mustInt(slice.GetInt()), t)
	ASSERT_EQ(int64(0x5a3fffac664223), mustInt(slice.GetSmallInt()), t)
}

func TestSliceInt8(t *testing.T) {
	slice := velocypack.Slice{0x27, 0x23, 0x42, 0x66, 0xac, 0xff, 0x3f, 0xfa, 0x6f}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(9), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(int64(0x6ffa3fffac664223), mustInt(slice.GetInt()), t)
	ASSERT_EQ(int64(0x6ffa3fffac664223), mustInt(slice.GetSmallInt()), t)
}

func TestSliceIntMax(t *testing.T) {
	t.Skip("TODO")
	/*	  Builder b;
	  b.add(Value(INT64_MAX));

	  Slice slice(b.slice());

		ASSERT_EQ(velocypack.Int, slice.Type(), t)
		ASSERT_TRUE(slice.IsInt(), t)
		ASSERT_EQ(velocypack.ValueLength(9), mustLength(slice.ByteSize()), t)

		ASSERT_EQ(int64(math.MaxInt64), mustInt(slice.GetInt()), t)
	*/
}

func TestSliceNegInt1(t *testing.T) {
	slice := velocypack.Slice{0x20, 0xa3}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(2), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(staticCastInt64(0xffffffffffffffa3), mustInt(slice.GetInt()), t)
}

func TestSliceNegInt2(t *testing.T) {
	slice := velocypack.Slice{0x21, 0x23, 0xe2}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(3), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(staticCastInt64(0xffffffffffffe223), mustInt(slice.GetInt()), t)
}

func TestSliceNegInt3(t *testing.T) {
	slice := velocypack.Slice{0x22, 0x23, 0x42, 0xd6}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(4), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(staticCastInt64(0xffffffffffd64223), mustInt(slice.GetInt()), t)
}

func TestSliceNegInt4(t *testing.T) {
	slice := velocypack.Slice{0x23, 0x23, 0x42, 0x66, 0xac}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(5), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(staticCastInt64(0xffffffffac664223), mustInt(slice.GetInt()), t)
}

func TestSliceNegInt5(t *testing.T) {
	slice := velocypack.Slice{0x24, 0x23, 0x42, 0x66, 0xac, 0xff}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(6), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(staticCastInt64(0xffffffffac664223), mustInt(slice.GetInt()), t)
}

func TestSliceNegInt6(t *testing.T) {
	slice := velocypack.Slice{0x25, 0x23, 0x42, 0x66, 0xac, 0xff, 0xef}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(7), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(staticCastInt64(0xffffefffac664223), mustInt(slice.GetInt()), t)
}

func TestSliceNegInt7(t *testing.T) {
	slice := velocypack.Slice{0x26, 0x23, 0x42, 0x66, 0xac, 0xff, 0xef, 0xfa}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(8), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(staticCastInt64(0xfffaefffac664223), mustInt(slice.GetInt()), t)
}

func TestSliceNegInt8(t *testing.T) {
	slice := velocypack.Slice{0x27, 0x23, 0x42, 0x66, 0xac, 0xff, 0xef, 0xfa, 0x8e}

	ASSERT_EQ(velocypack.Int, slice.Type(), t)
	ASSERT_TRUE(slice.IsInt(), t)
	ASSERT_EQ(velocypack.ValueLength(9), mustLength(slice.ByteSize()), t)

	ASSERT_EQ(staticCastInt64(0x8efaefffac664223), mustInt(slice.GetInt()), t)
}
