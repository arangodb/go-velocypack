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

func TestBuilderEmptyObject(t *testing.T) {
	var b velocypack.Builder
	b.OpenObject()
	b.Close()

	s := b.MustSlice()
	ASSERT_TRUE(s.IsObject(), t)
	ASSERT_EQ(velocypack.ValueLength(0), s.MustLength(), t)
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
