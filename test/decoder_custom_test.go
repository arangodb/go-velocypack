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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	velocypack "github.com/arangodb/go-velocypack"
)

/*
type CustomStruct1 struct {
	Field1 int
}
*/

func (cs *CustomStruct1) UnmarshalVPack(slice velocypack.Slice) error {
	s, err := slice.GetString()
	if err != nil {
		return err
	}
	if s != "Hello world" {
		return fmt.Errorf("Expected 'Hello world' got '%s'", s)
	}
	cs.Field1 = 42
	return nil
}

func TestDecoderCustomStruct1(t *testing.T) {
	input := &CustomStruct1{
		Field1: 999,
	}
	bytes, err := velocypack.Marshal(input)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)
	expected := CustomStruct1{
		Field1: 42,
	}

	var v CustomStruct1
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

/*
type CustomStruct2 struct {
	Field CustomStruct1
}
// CustomStruct2.Field is not using a custom unmarshaler since only *CustomStruct1 implements the Unmarshal interface.
*/

func TestDecoderCustomStruct2(t *testing.T) {
	input := CustomStruct2{
		Field: CustomStruct1{
			Field1: 999222,
		},
	}
	bytes, err := velocypack.Marshal(input)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)
	expected := input
	expected.Field.Field1 = 42

	var v CustomStruct2
	ASSERT_VELOCYPACK_EXCEPTION(func(error) bool { return true }, t)(velocypack.Unmarshal(s, &v))
}

/*
type CustomStruct3 struct {
	Field *CustomStruct1
}
*/

func TestDecoderCustomStruct3(t *testing.T) {
	input := CustomStruct3{
		Field: &CustomStruct1{
			Field1: 999222,
		},
	}
	bytes, err := velocypack.Marshal(input)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)
	expected := input
	expected.Field.Field1 = 42

	var v CustomStruct3
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

/*
type CustomText1 struct {
	I int
}
*/
func (ct *CustomText1) UnmarshalText(text []byte) error {
	if !strings.HasPrefix(string(text), "key") {
		return fmt.Errorf("Expected 'key' prefix, got '%s'", string(text))
	}
	i, err := strconv.Atoi(strings.TrimPrefix(string(text), "key"))
	if err != nil {
		return fmt.Errorf("Expected integer after 'key' prefix, got '%s' (err: %v)", strings.TrimPrefix(string(text), "key"), err)
	}
	ct.I = i
	return nil
}

func TestDecoderCustomText1(t *testing.T) {
	expected := map[CustomText1]bool{
		CustomText1{7}: true,
		CustomText1{2}: false,
	}
	bytes, err := velocypack.Marshal(expected)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)

	var v map[CustomText1]bool
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

func (cs *CustomJSONStruct1) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != "Hello JSON" {
		return fmt.Errorf("Expected 'Hello JSON' got '%s'", s)
	}
	cs.Field1 = 88
	return nil
}

func TestDecoderCustomJSONStruct1(t *testing.T) {
	input := &CustomJSONStruct1{
		Field1: 999,
	}
	bytes, err := velocypack.Marshal(input)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)
	expected := CustomJSONStruct1{
		Field1: 88,
	}

	var v CustomJSONStruct1
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}

func (cs *CustomJSONVPACKStruct1) UnmarshalVPack(slice velocypack.Slice) error {
	s, err := slice.GetString()
	if err != nil {
		return err
	}
	if s != "Hello VPACK, goodbye JSON" {
		return fmt.Errorf("Expected 'Hello VPACK, goodbye JSON' got '%s'", s)
	}
	cs.Field1 = 99
	return nil
}

func (cs *CustomJSONVPACKStruct1) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != "Hello JSON, goodbye VPACK" {
		return fmt.Errorf("Expected 'Hello JSON, goodbye VPACK' got '%s'", s)
	}
	cs.Field1 = 88
	return nil
}

func TestDecoderCustomJSONVPACKStruct1(t *testing.T) {
	// UnmarshalVPack is preferred over UnmarshalJSON
	input := &CustomJSONVPACKStruct1{
		Field1: 999,
	}
	bytes, err := velocypack.Marshal(input)
	ASSERT_NIL(err, t)
	s := velocypack.Slice(bytes)
	expected := CustomJSONVPACKStruct1{
		Field1: 99,
	}

	var v CustomJSONVPACKStruct1
	err = velocypack.Unmarshal(s, &v)
	ASSERT_NIL(err, t)
	ASSERT_EQ(v, expected, t)
}
