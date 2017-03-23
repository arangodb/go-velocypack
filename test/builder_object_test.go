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
	"encoding/binary"
	"math"
	"testing"

	velocypack "github.com/arangodb/go-velocypack"
)

func TestBuilderEmptyObject(t *testing.T) {
	var b velocypack.Builder
	b.OpenObject()
	b.Close()

	s := b.MustSlice()
	ASSERT_TRUE(s.IsObject(), t)
	ASSERT_EQ(velocypack.ValueLength(0), s.MustLength(), t)
}

func TestBuilderObjectEmpty(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewObjectValue())
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{0x0a}

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderObjectEmptyCompact(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewObjectValue(true))
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{0x0a}

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderObjectSorted(t *testing.T) {
	var b velocypack.Builder
	value := 2.3
	b.MustAddValue(velocypack.NewObjectValue())
	b.MustAddKeyValue("d", velocypack.NewUIntValue(1200))
	b.MustAddKeyValue("c", velocypack.NewDoubleValue(value))
	b.MustAddKeyValue("b", velocypack.NewStringValue("abc"))
	b.MustAddKeyValue("a", velocypack.NewBoolValue(true))
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{
		0x0b, 0x20, 0x04, 0x41, 0x64, 0x29, 0xb0, 0x04, // "d": uint(1200) =
		// 0x4b0
		0x41, 0x63, 0x1b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// "c": double(2.3)
		0x41, 0x62, 0x43, 0x61, 0x62, 0x63, // "b": "abc"
		0x41, 0x61, 0x1a, // "a": true
		0x19, 0x13, 0x08, 0x03}
	binary.LittleEndian.PutUint64(correctResult[11:], math.Float64bits(value))

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderObjectCompact(t *testing.T) {
	var b velocypack.Builder
	value := 2.3
	b.MustAddValue(velocypack.NewObjectValue(true))
	b.MustAddKeyValue("d", velocypack.NewUIntValue(1200))
	b.MustAddKeyValue("c", velocypack.NewDoubleValue(value))
	b.MustAddKeyValue("b", velocypack.NewStringValue("abc"))
	b.MustAddKeyValue("a", velocypack.NewBoolValue(true))
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{
		0x14, 0x1c, 0x41, 0x64, 0x29, 0xb0, 0x04, 0x41, 0x63, 0x1b,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // double
		0x41, 0x62, 0x43, 0x61, 0x62, 0x63, 0x41, 0x61, 0x1a, 0x04}
	binary.LittleEndian.PutUint64(correctResult[10:], math.Float64bits(value))

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderObjectValue1(t *testing.T) {
	var b velocypack.Builder
	u := uint64(77)
	b.OpenObject()
	b.AddKeyValue("test", velocypack.NewUIntValue(u))
	b.Close()

	s := b.MustSlice()
	ASSERT_TRUE(s.IsObject(), t)
	ASSERT_EQ(velocypack.ValueLength(1), s.MustLength(), t)
	ASSERT_EQ(u, s.MustGet("test").MustGetUInt(), t)
}

func TestBuilderObjectValue2(t *testing.T) {
	var b velocypack.Builder
	u := uint64(77)
	b.OpenObject()
	b.AddKeyValue("test", velocypack.NewUIntValue(u))
	b.AddKeyValue("soup", velocypack.NewUIntValue(u*2))
	b.Close()

	s := b.MustSlice()
	ASSERT_TRUE(s.IsObject(), t)
	ASSERT_EQ(velocypack.ValueLength(2), s.MustLength(), t)
	ASSERT_EQ(u, s.MustGet("test").MustGetUInt(), t)
	ASSERT_EQ(u*2, s.MustGet("soup").MustGetUInt(), t)
}

func TestBuilderAddObjectIteratorEmpty(t *testing.T) {
	var obj velocypack.Builder
	obj.OpenObject()
	obj.AddKeyValue("1-one", velocypack.NewIntValue(1))
	obj.AddKeyValue("2-two", velocypack.NewIntValue(2))
	obj.AddKeyValue("3-three", velocypack.NewIntValue(3))
	obj.Close()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	ASSERT_TRUE(b.IsClosed(), t)
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedOpenObject, t)(b.AddKeyValuesFromIterator(velocypack.MustNewObjectIterator(objSlice)))
	ASSERT_TRUE(b.IsClosed(), t)
}

func TestBuilderAddObjectIteratorKeyAlreadyWritten(t *testing.T) {
	var obj velocypack.Builder
	obj.OpenObject()
	obj.AddKeyValue("1-one", velocypack.NewIntValue(1))
	obj.AddKeyValue("2-two", velocypack.NewIntValue(2))
	obj.AddKeyValue("3-three", velocypack.NewIntValue(3))
	obj.Close()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	ASSERT_TRUE(b.IsClosed(), t)
	b.MustOpenObject()
	b.MustAddValue(velocypack.NewStringValue("foo"))
	ASSERT_FALSE(b.IsClosed(), t)
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderKeyAlreadyWritten, t)(b.AddKeyValuesFromIterator(velocypack.MustNewObjectIterator(objSlice)))
	ASSERT_FALSE(b.IsClosed(), t)
}

func TestBuilderAddObjectIteratorNonObject(t *testing.T) {
	var obj velocypack.Builder
	obj.OpenObject()
	obj.AddKeyValue("1-one", velocypack.NewIntValue(1))
	obj.AddKeyValue("2-two", velocypack.NewIntValue(2))
	obj.AddKeyValue("3-three", velocypack.NewIntValue(3))
	obj.Close()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	b.MustOpenArray()
	ASSERT_FALSE(b.IsClosed(), t)
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedOpenObject, t)(b.AddKeyValuesFromIterator(velocypack.MustNewObjectIterator(objSlice)))
	ASSERT_FALSE(b.IsClosed(), t)
}

func TestBuilderAddObjectIteratorTop(t *testing.T) {
	var obj velocypack.Builder
	obj.OpenObject()
	obj.AddKeyValue("1-one", velocypack.NewIntValue(1))
	obj.AddKeyValue("2-two", velocypack.NewIntValue(2))
	obj.AddKeyValue("3-three", velocypack.NewIntValue(3))
	obj.Close()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	b.MustOpenObject()
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustAddKeyValuesFromIterator(velocypack.MustNewObjectIterator(objSlice))
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose()
	result := b.MustSlice()
	ASSERT_TRUE(b.IsClosed(), t)

	ASSERT_EQ("{\"1-one\":1,\"2-two\":2,\"3-three\":3}", result.MustJSONString(), t)
}

func TestBuilderAddObjectIteratorReference(t *testing.T) {
	var obj velocypack.Builder
	obj.OpenObject()
	obj.AddKeyValue("1-one", velocypack.NewIntValue(1))
	obj.AddKeyValue("2-two", velocypack.NewIntValue(2))
	obj.AddKeyValue("3-three", velocypack.NewIntValue(3))
	obj.Close()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	b.MustOpenObject()
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustAdd(velocypack.MustNewObjectIterator(objSlice))
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose()
	result := b.MustSlice()
	ASSERT_TRUE(b.IsClosed(), t)

	ASSERT_EQ("{\"1-one\":1,\"2-two\":2,\"3-three\":3}", result.MustJSONString(), t)
}

func TestBuilderAddObjectIteratorSub(t *testing.T) {
	var obj velocypack.Builder
	obj.OpenObject()
	obj.AddKeyValue("1-one", velocypack.NewIntValue(1))
	obj.AddKeyValue("2-two", velocypack.NewIntValue(2))
	obj.AddKeyValue("3-three", velocypack.NewIntValue(3))
	obj.Close()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	b.MustOpenObject()
	b.MustAddKeyValue("1-something", velocypack.NewStringValue("tennis"))
	b.MustAddValue(velocypack.NewStringValue("2-values"))
	b.MustOpenObject()
	b.MustAdd(velocypack.MustNewObjectIterator(objSlice))
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose() // close one level
	b.MustAddKeyValue("3-bark", velocypack.NewStringValue("qux"))
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose()
	result := b.MustSlice()
	ASSERT_TRUE(b.IsClosed(), t)

	ASSERT_EQ("{\"1-something\":\"tennis\",\"2-values\":{\"1-one\":1,\"2-two\":2,\"3-three\":3},\"3-bark\":\"qux\"}", result.MustJSONString(), t)
}

func TestBuilderAddAndOpenObject(t *testing.T) {
	var b1 velocypack.Builder
	ASSERT_TRUE(b1.IsClosed(), t)
	b1.MustOpenObject()
	ASSERT_FALSE(b1.IsClosed(), t)
	b1.MustAddKeyValue("foo", velocypack.NewStringValue("bar"))
	b1.MustClose()
	ASSERT_TRUE(b1.IsClosed(), t)
	ASSERT_EQ(byte(0x14), b1.MustSlice()[0], t)
	ASSERT_EQ(velocypack.ValueLength(1), b1.MustSlice().MustLength(), t)

	var b2 velocypack.Builder
	ASSERT_TRUE(b2.IsClosed(), t)
	b2.MustOpenObject()
	ASSERT_FALSE(b2.IsClosed(), t)
	b2.MustAddKeyValue("foo", velocypack.NewStringValue("bar"))
	b2.MustClose()
	ASSERT_TRUE(b2.IsClosed(), t)
	ASSERT_EQ(byte(0x14), b2.MustSlice()[0], t)
	ASSERT_EQ(velocypack.ValueLength(1), b2.MustSlice().MustLength(), t)
}

func TestBuilderAddOnNonObject(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue())
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedOpenObject, t)(b.AddKeyValue("foo", velocypack.NewBoolValue(true)))
}
