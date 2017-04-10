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

func TestDecoderArrayEmpty(t *testing.T) {
	b := velocypack.Builder{}
	must(b.OpenArray())
	must(b.Close())
	s := mustSlice(b.Slice())

	var v []struct{}
	err := velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(0, len(v), t)
}

func TestDecoderArrayByteSlice(t *testing.T) {
	expected := []byte{1, 2, 3, 4, 5}
	b := velocypack.Builder{}
	must(b.AddValue(velocypack.NewBinaryValue(expected)))
	s := mustSlice(b.Slice())

	var v []byte
	err := velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

func TestDecoderArrayBoolSlice(t *testing.T) {
	expected := []bool{true, false, false, true}
	bytes, err := velocypack.Marshal(expected)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	var v []bool
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

func TestDecoderArrayIntSlice(t *testing.T) {
	expected := []int{1, 2, 3, -4, 5, 6, 100000}
	bytes, err := velocypack.Marshal(expected)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	var v []int
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

func TestDecoderArrayUIntSlice(t *testing.T) {
	expected := []uint{1, 2, 3, 4, 5, 6, 100000}
	bytes, err := velocypack.Marshal(expected)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	var v []uint
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

func TestDecoderArrayFloat32Slice(t *testing.T) {
	expected := []float32{0.0, -1.5, 66, 45}
	bytes, err := velocypack.Marshal(expected)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	var v []float32
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

func TestDecoderArrayFloat64Slice(t *testing.T) {
	expected := []float64{0.0, -1.5, 6.23, 45e+10}
	bytes, err := velocypack.Marshal(expected)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	var v []float64
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

/*
func TestDecoderArrayStructSlice(t *testing.T) {
	expected := []Struct1{
		Struct1{Field1: 1, field2: 2},
		Struct1{Field1: 10, field2: 200},
		Struct1{Field1: 100, field2: 200},
	}
	bytes, err := velocypack.Marshal(expected)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	var v []Struct1
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}


func TestDecoderArrayStructPtrSlice(t *testing.T) {
	bytes, err := velocypack.Marshal([]*Struct1{
		&Struct1{Field1: 1, field2: 2},
		nil,
		&Struct1{Field1: 10, field2: 200},
		&Struct1{Field1: 100, field2: 200},
		nil,
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	t.Log(s.String())
	ASSERT_EQ(s.Type(), velocypack.Array, t)
	ASSERT_TRUE(s.IsArray(), t)
	ASSERT_EQ(`[{"Field1":1},null,{"Field1":10},{"Field1":100},null]`, mustString(s.JSONString()), t)
}

func TestDecoderArrayNestedArray(t *testing.T) {
	bytes, err := velocypack.Marshal([][]Struct1{
		[]Struct1{Struct1{Field1: 1, field2: 2}, Struct1{Field1: 3, field2: 4}},
		[]Struct1{Struct1{Field1: 10, field2: 200}},
		[]Struct1{Struct1{Field1: 100, field2: 200}},
	})
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	t.Log(s.String())
	ASSERT_EQ(s.Type(), velocypack.Array, t)
	ASSERT_TRUE(s.IsArray(), t)
	ASSERT_EQ(`[[{"Field1":1},{"Field1":3}],[{"Field1":10}],[{"Field1":100}]]`, mustString(s.JSONString()), t)
}
*/
