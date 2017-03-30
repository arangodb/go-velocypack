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
	"encoding/hex"
	"testing"

	velocypack "github.com/arangodb/go-velocypack"
)

func TestBuilderAttributeTranslations(t *testing.T) {
	tx := velocypack.NewAttributeIDTranslator()
	tx.Add("foo", 1)
	tx.Add("bar", 2)
	tx.Add("baz", 3)
	tx.Add("bark", 4)
	tx.Add("mötör", 5)
	tx.Add("quetzalcoatl", 6)
	velocypack.AttributeTranslator = tx
	defer func() {
		velocypack.AttributeTranslator = nil
	}()

	var b velocypack.Builder
	b.AddValue(velocypack.NewObjectValue())
	b.AddKeyValue("foo", velocypack.NewBoolValue(true))
	b.AddKeyValue("bar", velocypack.NewBoolValue(false))
	b.AddKeyValue("baz", velocypack.NewIntValue(1))
	b.AddKeyValue("bart", velocypack.NewIntValue(2))
	b.AddKeyValue("bark", velocypack.NewIntValue(42))
	b.AddKeyValue("mötör", velocypack.NewIntValue(19))
	b.AddKeyValue("mötörhead", velocypack.NewIntValue(20))
	b.AddKeyValue("quetzal", velocypack.NewIntValue(21))
	must(b.Close())

	result := mustBytes(b.Bytes())

	correctResult := []byte{
		0x0b, 0x35, 0x08, 0x31, 0x1a, 0x32, 0x19, 0x33, 0x31, 0x44, 0x62,
		0x61, 0x72, 0x74, 0x32, 0x34, 0x20, 0x2a, 0x35, 0x20, 0x13, 0x4b,
		0x6d, 0xc3, 0xb6, 0x74, 0xc3, 0xb6, 0x72, 0x68, 0x65, 0x61, 0x64,
		0x20, 0x14, 0x47, 0x71, 0x75, 0x65, 0x74, 0x7a, 0x61, 0x6c, 0x20,
		0x15, 0x05, 0x0f, 0x09, 0x07, 0x03, 0x12, 0x15, 0x23}

	ASSERT_EQ(hex.EncodeToString(result), hex.EncodeToString(correctResult), t)

	s := mustSlice(b.Slice())
	ASSERT_TRUE(mustBool(s.HasKey("foo")), t)
	ASSERT_TRUE(mustBool(s.HasKey("bar")), t)
	ASSERT_TRUE(mustBool(s.HasKey("baz")), t)
	ASSERT_TRUE(mustBool(s.HasKey("bart")), t)
	ASSERT_TRUE(mustBool(s.HasKey("bark")), t)
	ASSERT_TRUE(mustBool(s.HasKey("mötör")), t)
	ASSERT_TRUE(mustBool(s.HasKey("mötörhead")), t)
	ASSERT_TRUE(mustBool(s.HasKey("quetzal")), t)
}

func TestBuilderAttributeTranslationsSorted(t *testing.T) {
	tx := velocypack.NewAttributeIDTranslator()
	tx.Add("foo", 1)
	tx.Add("bar", 2)
	tx.Add("baz", 3)
	tx.Add("bark", 4)
	tx.Add("mötör", 5)
	tx.Add("quetzalcoatl", 6)
	velocypack.AttributeTranslator = tx
	defer func() {
		velocypack.AttributeTranslator = nil
	}()

	var b velocypack.Builder
	b.AddValue(velocypack.NewObjectValue())
	b.AddKeyValue("foo", velocypack.NewBoolValue(true))
	b.AddKeyValue("bar", velocypack.NewBoolValue(false))
	b.AddKeyValue("baz", velocypack.NewIntValue(1))
	b.AddKeyValue("bart", velocypack.NewIntValue(2))
	b.AddKeyValue("bark", velocypack.NewIntValue(42))
	b.AddKeyValue("mötör", velocypack.NewIntValue(19))
	b.AddKeyValue("mötörhead", velocypack.NewIntValue(20))
	b.AddKeyValue("quetzal", velocypack.NewIntValue(21))
	must(b.Close())

	result := mustBytes(b.Bytes())

	correctResult := []byte{
		0x0b, 0x35, 0x08, 0x31, 0x1a, 0x32, 0x19, 0x33, 0x31, 0x44, 0x62,
		0x61, 0x72, 0x74, 0x32, 0x34, 0x20, 0x2a, 0x35, 0x20, 0x13, 0x4b,
		0x6d, 0xc3, 0xb6, 0x74, 0xc3, 0xb6, 0x72, 0x68, 0x65, 0x61, 0x64,
		0x20, 0x14, 0x47, 0x71, 0x75, 0x65, 0x74, 0x7a, 0x61, 0x6c, 0x20,
		0x15, 0x05, 0x0f, 0x09, 0x07, 0x03, 0x12, 0x15, 0x23}

	ASSERT_EQ(hex.EncodeToString(result), hex.EncodeToString(correctResult), t)

	s := mustSlice(b.Slice())
	ASSERT_TRUE(mustBool(s.HasKey("foo")), t)
	ASSERT_TRUE(mustBool(s.HasKey("bar")), t)
	ASSERT_TRUE(mustBool(s.HasKey("baz")), t)
	ASSERT_TRUE(mustBool(s.HasKey("bart")), t)
	ASSERT_TRUE(mustBool(s.HasKey("bark")), t)
	ASSERT_TRUE(mustBool(s.HasKey("mötör")), t)
	ASSERT_TRUE(mustBool(s.HasKey("mötörhead")), t)
	ASSERT_TRUE(mustBool(s.HasKey("quetzal")), t)
}
