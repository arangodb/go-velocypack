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

func TestBuilderBytesWithOpenObject(t *testing.T) {
	var b velocypack.Builder
	ASSERT_EQ(0, len(b.MustBytes()), t)
	b.MustOpenObject()
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNotSealed, t)(b.Bytes())
	b.MustClose()
	ASSERT_EQ(1, len(b.MustBytes()), t)
}

func TestBuilderSliceWithOpenObject(t *testing.T) {
	var b velocypack.Builder
	ASSERT_EQ(0, len(b.MustSlice()), t)
	b.MustOpenObject()
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNotSealed, t)(b.Slice())
	b.MustClose()
	ASSERT_EQ(1, len(b.MustSlice()), t)
}

func TestBuilderSizeWithOpenObject(t *testing.T) {
	var b velocypack.Builder
	ASSERT_EQ(velocypack.ValueLength(0), b.MustSize(), t)
	b.MustOpenObject()
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNotSealed, t)(b.Size())
	b.MustClose()
	ASSERT_EQ(velocypack.ValueLength(1), b.MustSize(), t)
}

func TestBuilderIsEmpty(t *testing.T) {
	var b velocypack.Builder
	ASSERT_TRUE(b.IsEmpty(), t)
	b.MustOpenObject()
	ASSERT_FALSE(b.IsEmpty(), t)
}

func TestBuilderIsClosedMixed(t *testing.T) {
	var b velocypack.Builder
	ASSERT_TRUE(b.IsClosed(), t)
	b.AddValue(velocypack.NewNullValue())
	ASSERT_TRUE(b.IsClosed(), t)
	b.AddValue(velocypack.NewBoolValue(true))
	ASSERT_TRUE(b.IsClosed(), t)

	b.AddValue(velocypack.NewArrayValue())
	ASSERT_FALSE(b.IsClosed(), t)

	b.AddValue(velocypack.NewBoolValue(true))
	ASSERT_FALSE(b.IsClosed(), t)
	b.AddValue(velocypack.NewBoolValue(true))
	ASSERT_FALSE(b.IsClosed(), t)

	b.MustClose()
	ASSERT_TRUE(b.IsClosed(), t)

	b.AddValue(velocypack.NewObjectValue())
	ASSERT_FALSE(b.IsClosed(), t)

	b.AddKeyValue("foo", velocypack.NewBoolValue(true))
	ASSERT_FALSE(b.IsClosed(), t)

	b.AddKeyValue("bar", velocypack.NewBoolValue(true))
	ASSERT_FALSE(b.IsClosed(), t)

	b.AddKeyValue("baz", velocypack.NewArrayValue())
	ASSERT_FALSE(b.IsClosed(), t)

	b.MustClose()
	ASSERT_FALSE(b.IsClosed(), t)

	b.MustClose()
	ASSERT_TRUE(b.IsClosed(), t)
}

func TestBuilderIsClosedObject(t *testing.T) {
	var b velocypack.Builder
	ASSERT_TRUE(b.IsClosed(), t)
	b.MustAddValue(velocypack.NewObjectValue())
	ASSERT_FALSE(b.IsClosed(), t)

	b.MustAddKeyValue("foo", velocypack.NewBoolValue(true))
	ASSERT_FALSE(b.IsClosed(), t)

	b.MustAddKeyValue("bar", velocypack.NewBoolValue(true))
	ASSERT_FALSE(b.IsClosed(), t)

	b.MustAddKeyValue("baz", velocypack.NewObjectValue())
	ASSERT_FALSE(b.IsClosed(), t)

	b.MustClose()
	ASSERT_FALSE(b.IsClosed(), t)

	b.MustClose()
	ASSERT_TRUE(b.IsClosed(), t)
}

func TestBuilderCloseClosed(t *testing.T) {
	var b velocypack.Builder
	ASSERT_TRUE(b.IsClosed(), t)
	b.MustAddValue(velocypack.NewObjectValue())
	ASSERT_FALSE(b.IsClosed(), t)
	b.MustClose()
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedOpenCompound, t)(b.Close())
}

func TestBuilderRemoveLastNonObject(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewBoolValue(true))
	b.MustAddValue(velocypack.NewBoolValue(false))
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedOpenCompound, t)(b.RemoveLast())
}

func TestBuilderRemoveLastSealed(t *testing.T) {
	var b velocypack.Builder
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedOpenCompound, t)(b.RemoveLast())
}

func TestBuilderRemoveLastEmptyObject(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewObjectValue())
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedSubValue, t)(b.RemoveLast())
}

func TestBuilderRemoveLastObjectInvalid(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewObjectValue())
	b.MustAddKeyValue("foo", velocypack.NewBoolValue(true))
	b.MustRemoveLast()
	ASSERT_VELOCYPACK_EXCEPTION(velocypack.IsBuilderNeedSubValue, t)(b.RemoveLast())
}

func TestBuilderRemoveLastObject(t *testing.T) {
	var b velocypack.Builder
	b.MustAddValue(velocypack.NewObjectValue())
	b.MustAddKeyValue("foo", velocypack.NewBoolValue(true))
	b.MustAddKeyValue("bar", velocypack.NewBoolValue(false))

	b.MustRemoveLast()
	b.MustClose()

	s := b.MustSlice()
	ASSERT_TRUE(s.IsObject(), t)
	ASSERT_EQ(velocypack.ValueLength(1), s.MustLength(), t)
	ASSERT_TRUE(s.MustHasKey("foo"), t)
	ASSERT_TRUE(s.MustGet("foo").MustGetBool(), t)
	ASSERT_FALSE(s.MustHasKey("bar"), t)
}
