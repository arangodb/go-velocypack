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

func TestSliceStringNoString(t *testing.T) {
	slice := velocypack.Slice{}

	ASSERT_FALSE(slice.IsString(), t)
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsInvalidType, t)(slice.GetString())
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsInvalidType, t)(slice.GetStringLength())
}

func TestSliceStringEmpty(t *testing.T) {
	slice := velocypack.Slice{0x40}

	ASSERT_EQ(velocypack.String, slice.Type(), t)
	ASSERT_TRUE(slice.IsString(), t)
	ASSERT_EQ(velocypack.ValueLength(1), slice.MustByteSize(), t)
	ASSERT_EQ("", slice.MustGetString(), t)
	ASSERT_EQ(velocypack.ValueLength(0), slice.MustGetStringLength(), t)
}

func TestSliceStringLengths(t *testing.T) {
	t.Skip("TODO")
	/*
			Builder builder;

		  for (size_t i = 0; i < 255; ++i) {
		    builder.clear();

		    std::string temp;
		    for (size_t j = 0; j < i; ++j) {
		      temp.push_back('x');
		    }

		    builder.add(Value(temp));

		    Slice slice = builder.slice();

		    ASSERT_TRUE(slice.isString());
		    ASSERT_EQ(ValueType::String, slice.type());

		    ASSERT_EQ(i, slice.getStringLength());

		    if (i <= 126) {
		      ASSERT_EQ(i + 1, slice.byteSize());
		    } else {
		      ASSERT_EQ(i + 9, slice.byteSize());
		    }
		  }
	*/
}

func TestSliceString1(t *testing.T) {
	value := "foobar"
	slice := velocypack.Slice(append([]byte{byte(0x40 + len(value))}, value...))

	ASSERT_EQ(velocypack.String, slice.Type(), t)
	ASSERT_TRUE(slice.IsString(), t)
	ASSERT_EQ(velocypack.ValueLength(7), slice.MustByteSize(), t)
	ASSERT_EQ(value, slice.MustGetString(), t)
	ASSERT_EQ(velocypack.ValueLength(len(value)), slice.MustGetStringLength(), t)
}

func TestSliceString2(t *testing.T) {
	slice := velocypack.Slice{0x48, '1', '2', '3', 'f', '\r', '\t', '\n', 'x'}

	ASSERT_EQ(velocypack.String, slice.Type(), t)
	ASSERT_TRUE(slice.IsString(), t)
	ASSERT_EQ(velocypack.ValueLength(9), slice.MustByteSize(), t)
	ASSERT_EQ("123f\r\t\nx", slice.MustGetString(), t)
	ASSERT_EQ(velocypack.ValueLength(8), slice.MustGetStringLength(), t)
}

func TestSliceStringNullBytes(t *testing.T) {
	slice := velocypack.Slice{0x48, 0, '1', '2', 0, '3', '4', 0, 'x'}

	ASSERT_EQ(velocypack.String, slice.Type(), t)
	ASSERT_TRUE(slice.IsString(), t)
	ASSERT_EQ(velocypack.ValueLength(9), slice.MustByteSize(), t)
	ASSERT_EQ("\x0012\x0034\x00x", slice.MustGetString(), t)
	ASSERT_EQ(velocypack.ValueLength(8), slice.MustGetStringLength(), t)
}

func TestSliceStringLong(t *testing.T) {
	slice := velocypack.Slice{0xbf, 6, 0, 0, 0, 0, 0, 0, 0, 'f', 'o', 'o', 'b', 'a', 'r'}

	ASSERT_EQ(velocypack.String, slice.Type(), t)
	ASSERT_TRUE(slice.IsString(), t)
	ASSERT_EQ(velocypack.ValueLength(15), slice.MustByteSize(), t)
	ASSERT_EQ("foobar", slice.MustGetString(), t)
	ASSERT_EQ(velocypack.ValueLength(6), slice.MustGetStringLength(), t)
}

func TestSliceStringToStringNull(t *testing.T) {
	slice := velocypack.NullSlice()

	ASSERT_EQ("null", slice.MustJSONString(), t)
}
