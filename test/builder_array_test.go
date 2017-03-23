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

func TestBuilderEmptyArray(t *testing.T) {
	var b velocypack.Builder
	b.OpenArray()
	b.Close()

	s := b.MustSlice()
	ASSERT_TRUE(s.IsArray(), t)
	ASSERT_EQ(velocypack.ValueLength(0), s.MustLength(), t)
}

func TestBuilderArrayEmpty(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue())
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{0x01}

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderArraySingleEntry(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue())
	b.MustAddValue(velocypack.NewIntValue(1))
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{0x02, 0x03, 0x31}

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderArraySingleEntryLong(t *testing.T) {
	value := "ngdddddljjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjsdddffffffffffffmmmmmmmmmmmmmmmsf" +
		"dlllllllllllllllllllllllllllllllllllllllllllllllllrjjjjjjsdddddddddddddd" +
		"ddddhhhhhhkkkkkkkksssssssssssssssssssssssssssssssssdddddddddddddddddkkkk" +
		"kkkkkkkkksddddddddddddssssssssssfvvvvvvvvvvvvvvvvvvvvvvvvvvvfvgfff"
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue())
	b.MustAddValue(velocypack.NewStringValue(value))
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{
		0x03, 0x2c, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xbf, 0x1a, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x6e, 0x67, 0x64, 0x64, 0x64, 0x64,
		0x64, 0x6c, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a,
		0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a,
		0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x6a, 0x73, 0x64, 0x64,
		0x64, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
		0x66, 0x6d, 0x6d, 0x6d, 0x6d, 0x6d, 0x6d, 0x6d, 0x6d, 0x6d, 0x6d, 0x6d,
		0x6d, 0x6d, 0x6d, 0x6d, 0x73, 0x66, 0x64, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c,
		0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c,
		0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c,
		0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c,
		0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x6c, 0x72, 0x6a, 0x6a, 0x6a,
		0x6a, 0x6a, 0x6a, 0x73, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64,
		0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x68, 0x68,
		0x68, 0x68, 0x68, 0x68, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b,
		0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73,
		0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73,
		0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x64, 0x64, 0x64,
		0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64,
		0x64, 0x64, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b, 0x6b,
		0x6b, 0x6b, 0x6b, 0x73, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64,
		0x64, 0x64, 0x64, 0x64, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73,
		0x73, 0x73, 0x66, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76,
		0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x76,
		0x76, 0x76, 0x76, 0x76, 0x76, 0x76, 0x66, 0x76, 0x67, 0x66, 0x66, 0x66}

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderArraySameSizeEntries(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue())
	b.MustAddValue(velocypack.NewUIntValue(1))
	b.MustAddValue(velocypack.NewUIntValue(2))
	b.MustAddValue(velocypack.NewUIntValue(3))
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{0x02, 0x05, 0x31, 0x32, 0x33}

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderArraySomeEntries(t *testing.T) {
	var b velocypack.Builder
	value := 2.3
	b.MustAddValue(velocypack.NewArrayValue())
	b.MustAddValue(velocypack.NewUIntValue(1200))
	b.MustAddValue(velocypack.NewDoubleValue(value))
	b.MustAddValue(velocypack.NewStringValue("abc"))
	b.MustAddValue(velocypack.NewBoolValue(true))
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{
		0x06, 0x18, 0x04, 0x29, 0xb0, 0x04, // uint(1200) = 0x4b0
		0x1b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // double(2.3)
		0x43, 0x61, 0x62, 0x63, 0x1a, 0x03, 0x06, 0x0f, 0x13}
	binary.LittleEndian.PutUint64(correctResult[7:], math.Float64bits(value))

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderArrayCompact(t *testing.T) {
	var b velocypack.Builder
	value := 2.3
	b.MustAddValue(velocypack.NewArrayValue(true))
	b.MustAddValue(velocypack.NewUIntValue(1200))
	b.MustAddValue(velocypack.NewDoubleValue(value))
	b.MustAddValue(velocypack.NewStringValue("abc"))
	b.MustAddValue(velocypack.NewBoolValue(true))
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	correctResult := []byte{
		0x13, 0x14, 0x29, 0xb0, 0x04, 0x1b,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, // double
		0x43, 0x61, 0x62, 0x63, 0x1a, 0x04}
	binary.LittleEndian.PutUint64(correctResult[6:], math.Float64bits(value))

	ASSERT_EQ(velocypack.ValueLength(len(correctResult)), l, t)
	ASSERT_EQ(result, correctResult, t)
}

func TestBuilderArrayCompactBytesizeBelowThreshold(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue(true))
	for i := uint64(0); i < 124; i++ {
		b.MustAddValue(velocypack.NewUIntValue(i % 10))
	}
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	ASSERT_EQ(velocypack.ValueLength(127), l, t)
	ASSERT_EQ(byte(0x13), result[0], t)
	ASSERT_EQ(byte(0x7f), result[1], t)
	for i := uint64(0); i < 124; i++ {
		ASSERT_EQ(byte(0x30+(i%10)), result[2+i], t)
	}
	ASSERT_EQ(byte(0x7c), result[126], t)
}

func TestBuilderArrayCompactBytesizeAboveThreshold(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue(true))
	for i := uint64(0); i < 125; i++ {
		b.MustAddValue(velocypack.NewUIntValue(i % 10))
	}
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	ASSERT_EQ(velocypack.ValueLength(129), l, t)
	ASSERT_EQ(byte(0x13), result[0], t)
	ASSERT_EQ(byte(0x81), result[1], t)
	ASSERT_EQ(byte(0x01), result[2], t)
	for i := uint64(0); i < 125; i++ {
		ASSERT_EQ(byte(0x30+(i%10)), result[3+i], t)
	}
	ASSERT_EQ(byte(0x7d), result[128], t)
}

func TestBuilderArrayCompactLengthBelowThreshold(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue(true))
	for i := uint64(0); i < 127; i++ {
		b.MustAddValue(velocypack.NewStringValue("aaa"))
	}
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	ASSERT_EQ(velocypack.ValueLength(512), l, t)
	ASSERT_EQ(byte(0x13), result[0], t)
	ASSERT_EQ(byte(0x80), result[1], t)
	ASSERT_EQ(byte(0x04), result[2], t)
	for i := uint64(0); i < 127; i++ {
		ASSERT_EQ(byte(0x43), result[3+i*4], t)
	}
	ASSERT_EQ(byte(0x7f), result[511], t)
}

func TestBuilderArrayCompactLengthAboveThreshold(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewArrayValue(true))
	for i := uint64(0); i < 128; i++ {
		b.MustAddValue(velocypack.NewStringValue("aaa"))
	}
	b.MustClose()
	l := b.MustSize()
	result := b.MustBytes()

	ASSERT_EQ(velocypack.ValueLength(517), l, t)
	ASSERT_EQ(byte(0x13), result[0], t)
	ASSERT_EQ(byte(0x85), result[1], t)
	ASSERT_EQ(byte(0x04), result[2], t)
	for i := uint64(0); i < 128; i++ {
		ASSERT_EQ(byte(0x43), result[3+i*4], t)
	}
	ASSERT_EQ(byte(0x01), result[515], t)
	ASSERT_EQ(byte(0x80), result[516], t)
}

func TestBuilderAddObjectInArray(t *testing.T) {
	var b velocypack.Builder
	b.OpenArray()
	b.OpenObject()
	b.Close()
	b.Close()

	s := b.MustSlice()
	ASSERT_TRUE(s.IsArray(), t)
	ASSERT_EQ(velocypack.ValueLength(1), s.MustLength(), t)
	ss := s.MustAt(0)
	ASSERT_TRUE(ss.IsObject(), t)
	ASSERT_EQ(velocypack.ValueLength(0), ss.MustLength(), t)
}

func TestBuilderAddArrayIteratorEmpty(t *testing.T) {
	var obj velocypack.Builder
	obj.MustOpenArray()
	obj.MustAddValue(velocypack.NewIntValue(1))
	obj.MustAddValue(velocypack.NewIntValue(2))
	obj.MustAddValue(velocypack.NewIntValue(3))
	obj.MustClose()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	ASSERT_TRUE(b.IsClosed(), t)
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedOpenArray, t)(b.AddValuesFromIterator(velocypack.MustNewArrayIterator(objSlice)))
	ASSERT_TRUE(b.IsClosed(), t)
}

func TestBuilderAddArrayIteratorNonArray(t *testing.T) {
	var obj velocypack.Builder
	obj.MustOpenArray()
	obj.MustAddValue(velocypack.NewIntValue(1))
	obj.MustAddValue(velocypack.NewIntValue(2))
	obj.MustAddValue(velocypack.NewIntValue(3))
	obj.MustClose()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	b.MustOpenObject()
	ASSERT_FALSE(b.IsClosed(), t)
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedOpenArray, t)(b.AddValuesFromIterator(velocypack.MustNewArrayIterator(objSlice)))
	ASSERT_FALSE(b.IsClosed(), t)
}

func TestBuilderAddArrayIteratorTop(t *testing.T) {
	var obj velocypack.Builder
	obj.MustOpenArray()
	obj.MustAddValue(velocypack.NewIntValue(1))
	obj.MustAddValue(velocypack.NewIntValue(2))
	obj.MustAddValue(velocypack.NewIntValue(3))
	obj.MustClose()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	b.MustOpenArray()
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustAddValuesFromIterator(velocypack.MustNewArrayIterator(objSlice))
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose()
	result := b.MustSlice()

	ASSERT_EQ("[1,2,3]", result.MustJSONString(), t)
}

func TestBuilderAddArrayIteratorReference(t *testing.T) {
	var obj velocypack.Builder
	obj.MustOpenArray()
	obj.MustAddValue(velocypack.NewIntValue(1))
	obj.MustAddValue(velocypack.NewIntValue(2))
	obj.MustAddValue(velocypack.NewIntValue(3))
	obj.MustClose()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	b.MustOpenArray()
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustAdd(velocypack.MustNewArrayIterator(objSlice))
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose()
	result := b.MustSlice()

	ASSERT_EQ("[1,2,3]", result.MustJSONString(), t)
}

func TestBuilderAddArrayIteratorSub(t *testing.T) {
	var obj velocypack.Builder
	obj.MustOpenArray()
	obj.MustAddValue(velocypack.NewIntValue(1))
	obj.MustAddValue(velocypack.NewIntValue(2))
	obj.MustAddValue(velocypack.NewIntValue(3))
	obj.MustClose()
	objSlice := obj.MustSlice()

	var b velocypack.Builder
	b.MustOpenArray()
	b.MustAddValue(velocypack.NewStringValue("tennis"))
	b.MustOpenArray()
	b.MustAdd(velocypack.MustNewArrayIterator(objSlice))
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose() // close one level
	b.MustAddValue(velocypack.NewStringValue("qux"))
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose()
	result := b.MustSlice()
	ASSERT_TRUE(b.IsClosed(), t)

	ASSERT_EQ("[\"tennis\",[1,2,3],\"qux\"]", result.MustJSONString(), t)
}

func TestBuilderAddAndOpenArray(t *testing.T) {
	var b1 velocypack.Builder
	ASSERT_TRUE(b1.IsClosed(), t)
	b1.MustOpenArray()
	ASSERT_FALSE(b1.IsClosed(), t)
	b1.MustAddValue(velocypack.NewStringValue("bar"))
	b1.MustClose()
	ASSERT_TRUE(b1.IsClosed(), t)
	ASSERT_EQ(byte(0x02), b1.MustSlice()[0], t)

	var b2 velocypack.Builder
	ASSERT_TRUE(b2.IsClosed(), t)
	b2.MustOpenArray()
	ASSERT_FALSE(b2.IsClosed(), t)
	b2.MustAddValue(velocypack.NewStringValue("bar"))
	b2.MustClose()
	ASSERT_TRUE(b2.IsClosed(), t)
	ASSERT_EQ(byte(0x02), b2.MustSlice()[0], t)
}

func TestBuilderAddOnNonArray(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewObjectValue())
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderKeyMustBeString, t)(b.AddValue(velocypack.NewBoolValue(true)))
}
