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

// Slice provides read only access to a VPack value
type Slice []byte

// Type returns the vpack type of the slice
func (s Slice) Type() ValueType {
	return typeMap[s[0]]
}

// IsType returns true when the vpack type of the slice is equal to the given type.
// Returns false otherwise.
func (s Slice) IsType(t ValueType) bool {
	return typeMap[s[0]] == t
}

// check if slice is a None object
func (s Slice) isNone() bool { return s.IsType(None) }

// check if slice is an Illegal object
func (s Slice) isIllegal() bool { return s.IsType(Illegal) }

// check if slice is a Null object
func (s Slice) isNull() bool { return s.IsType(Null) }

// check if slice is a Bool object
func (s Slice) isBool() bool { return s.IsType(Bool) }

// check if slice is the Boolean value true
func (s Slice) isTrue() bool { return s.head() == 0x1a }

// check if slice is the Boolean value false
func (s Slice) isFalse() bool { return s.head() == 0x19 }

// check if slice is an Array object
func (s Slice) isArray() bool { return s.IsType(Array) }

// check if slice is an Object object
func (s Slice) isObject() bool { return s.IsType(Object) }

// check if slice is a Double object
func (s Slice) isDouble() bool { return s.IsType(Double) }

// check if slice is a UTCDate object
func (s Slice) isUTCDate() bool { return s.IsType(UTCDate) }

// check if slice is an External object
func (s Slice) isExternal() bool { return s.IsType(External) }

// check if slice is a MinKey object
func (s Slice) isMinKey() bool { return s.IsType(MinKey) }

// check if slice is a MaxKey object
func (s Slice) isMaxKey() bool { return s.IsType(MaxKey) }

// check if slice is an Int object
func (s Slice) isInt() bool { return s.IsType(Int) }

// check if slice is a UInt object
func (s Slice) isUInt() bool { return s.IsType(UInt) }

// check if slice is a SmallInt object
func (s Slice) isSmallInt() bool { return s.IsType(SmallInt) }

// check if slice is a String object
func (s Slice) isString() bool { return s.IsType(String) }

// check if slice is a Binary object
func (s Slice) isBinary() bool { return s.IsType(Binary) }

// check if slice is a BCD
func (s Slice) isBCD() bool { return s.IsType(BCD) }

// check if slice is a Custom type
func (s Slice) isCustom() bool { return s.IsType(Custom) }

// check if a slice is any number type
func (s Slice) isInteger() bool {
	return (s.isInt() || s.isUInt() || s.isSmallInt())
}

// check if slice is any Number-type object
func (s Slice) isNumber() bool { return s.isInteger() || s.isDouble() }

func (s Slice) isSorted() bool {
	h := s.head()
	return (h >= 0x0b && h <= 0x0e)
}
