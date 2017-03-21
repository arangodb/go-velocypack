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

package velocypack

import "time"

// Value is a helper structure used to build VPack structures.
// It holds a single data value with a specific type.
type Value struct {
	vt        ValueType
	data      interface{}
	unindexed bool
}

// NewBoolValue creates a new Value of type Bool with given value.
func NewBoolValue(value bool) Value {
	return Value{Bool, value, false}
}

// NewIntValue creates a new Value of type Int with given value.
func NewIntValue(value int64) Value {
	return Value{Int, value, false}
}

// NewUIntValue creates a new Value of type UInt with given value.
func NewUIntValue(value uint64) Value {
	return Value{UInt, value, false}
}

// NewDoubleValue creates a new Value of type Double with given value.
func NewDoubleValue(value float64) Value {
	return Value{Double, value, false}
}

// NewStringValue creates a new Value of type String with given value.
func NewStringValue(value string) Value {
	return Value{String, value, false}
}

// NewBinaryValue creates a new Value of type Binary with given value.
func NewBinaryValue(value []byte) Value {
	return Value{Binary, value, false}
}

// NewUTCDateValue creates a new Value of type UTCDate with given value.
func NewUTCDateValue(value time.Time) Value {
	return Value{UTCDate, value, false}
}

func (v Value) boolValue() bool {
	return v.data.(bool)
}

func (v Value) intValue() int64 {
	return v.data.(int64)
}

func (v Value) uintValue() uint64 {
	return v.data.(uint64)
}

func (v Value) doubleValue() float64 {
	return v.data.(float64)
}

func (v Value) stringValue() string {
	return v.data.(string)
}

func (v Value) binaryValue() []byte {
	return v.data.([]byte)
}

func (v Value) utcDateValue() int64 {
	time := v.data.(time.Time)
	return time.Unix()
}
