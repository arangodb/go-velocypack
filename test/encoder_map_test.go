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

func TestEncoderMapEmpty(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_TRUE(s.IsEmptyObject(), t)
}

func TestEncoderMapOneField(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]string{
		"Name": "Max",
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"Name":"Max"}`, mustString(s.JSONString()), t)
}

func TestEncoderMapMultipleFields(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name": "Max",
		"A":    true,
		"D":    123.456,
		"I":    789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Max"}`, mustString(s.JSONString()), t)
}

func TestEncoderMapMultipleFieldsEmpty(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name": "",
		"A":    false,
		"D":    0.0,
		"I":    0,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":false,"D":0,"I":0,"Name":""}`, mustString(s.JSONString()), t)
}

func TestEncoderMapNestedStruct(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name": "Jan",
		"Nested": map[string]interface{}{
			"Foo": 999,
		},
		"A": true,
		"D": 123.456,
		"I": 789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Jan","Nested":{"Foo":999}}`, mustString(s.JSONString()), t)
}

func TestEncoderMapNestedStructs(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name": "Jan",
		"Nested": map[string]interface{}{
			"Foo": 999,
			"Nested": map[string]bool{
				"Foo": true,
			},
		},
		"A": true,
		"D": 123.456,
		"I": 789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Jan","Nested":{"Foo":999,"Nested":{"Foo":true}}}`, mustString(s.JSONString()), t)
}

func TestEncoderMapNestedStructPtr(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name": "Jan",
		"Nested": &struct {
			Foo int
		}{
			Foo: 999,
		},
		"A": true,
		"D": 123.456,
		"I": 789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Jan","Nested":{"Foo":999}}`, mustString(s.JSONString()), t)
}

func TestEncoderMapNestedStructPtrNil(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name":   "Jan",
		"Nested": nil,
		"A":      true,
		"D":      123.456,
		"I":      789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Jan","Nested":null}`, mustString(s.JSONString()), t)
}

func TestEncoderMapNestedByteSlice(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name":   "Jan",
		"Nested": []byte{1, 2, 3, 4, 5, 6},
		"A":      true,
		"D":      123.456,
		"I":      789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Jan","Nested":"(non-representable type Binary)"}`, mustString(s.JSONString(velocypack.DumperOptions{UnsupportedTypeBehavior: velocypack.ConvertUnsupportedType})), t)
}

func TestEncoderMapNestedIntSlice(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name":   "Jan",
		"Nested": []int{1, 2, 3, 4, 5},
		"A":      true,
		"D":      123.456,
		"I":      789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Jan","Nested":[1,2,3,4,5]}`, mustString(s.JSONString()), t)
}

func TestEncoderMapNestedStringSlice(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name":   "Jan",
		"Nested": []string{"Aap", "Noot"},
		"A":      true,
		"D":      123.456,
		"I":      789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Jan","Nested":["Aap","Noot"]}`, mustString(s.JSONString()), t)
}

func TestEncoderMapNestedStringSliceEmpty(t *testing.T) {
	bytes, err := velocypack.Marshal(map[string]interface{}{
		"Name":   "Jan",
		"Nested": []string{},
		"A":      true,
		"D":      123.456,
		"I":      789,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	ASSERT_EQ(s.Type(), velocypack.Object, t)
	ASSERT_FALSE(s.IsEmptyObject(), t)
	ASSERT_EQ(`{"A":true,"D":123.456,"I":789,"Name":"Jan","Nested":[]}`, mustString(s.JSONString()), t)
}
