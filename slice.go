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

import "encoding/hex"

// Slice provides read only access to a VPack value
type Slice []byte

// SliceFromBytes creates a Slice by casting the given byte slice to a Slice.
func SliceFromBytes(v []byte) Slice {
	return Slice(v)
}

// SliceFromHex creates a Slice by decoding the given hex string into a Slice.
// If decoding fails, nil is returned.
func SliceFromHex(v string) Slice {
	if bytes, err := hex.DecodeString(v); err != nil {
		return nil
	} else {
		return Slice(bytes)
	}
}

// String returns a HEX representation of the slice.
func (s Slice) String() string {
	return hex.EncodeToString(s)
}

// Type returns the vpack type of the slice
func (s Slice) Type() ValueType {
	return typeMap[s[0]]
}

// IsType returns true when the vpack type of the slice is equal to the given type.
// Returns false otherwise.
func (s Slice) IsType(t ValueType) bool {
	return typeMap[s[0]] == t
}

// IsNone returns true if slice is a None object
func (s Slice) IsNone() bool { return s.IsType(None) }

// IsIllegal returns true if slice is an Illegal object
func (s Slice) IsIllegal() bool { return s.IsType(Illegal) }

// IsNull returns true if slice is a Null object
func (s Slice) IsNull() bool { return s.IsType(Null) }

// IsBool returns true if slice is a Bool object
func (s Slice) IsBool() bool { return s.IsType(Bool) }

// IsTrue returns true if slice is the Boolean value true
func (s Slice) IsTrue() bool { return s[0] == 0x1a }

// IsFalse returns true if slice is the Boolean value false
func (s Slice) IsFalse() bool { return s[0] == 0x19 }

// IsArray returns true if slice is an Array object
func (s Slice) IsArray() bool { return s.IsType(Array) }

// IsObject returns true if slice is an Object object
func (s Slice) IsObject() bool { return s.IsType(Object) }

// IsDouble returns true if slice is a Double object
func (s Slice) IsDouble() bool { return s.IsType(Double) }

// IsUTCDate returns true if slice is a UTCDate object
func (s Slice) IsUTCDate() bool { return s.IsType(UTCDate) }

// IsExternal returns true if slice is an External object
func (s Slice) IsExternal() bool { return s.IsType(External) }

// IsMinKey returns true if slice is a MinKey object
func (s Slice) IsMinKey() bool { return s.IsType(MinKey) }

// IsMaxKey returns true if slice is a MaxKey object
func (s Slice) IsMaxKey() bool { return s.IsType(MaxKey) }

// IsInt returns true if slice is an Int object
func (s Slice) IsInt() bool { return s.IsType(Int) }

// IsUInt returns true if slice is a UInt object
func (s Slice) IsUInt() bool { return s.IsType(UInt) }

// IsSmallInt returns true if slice is a SmallInt object
func (s Slice) IsSmallInt() bool { return s.IsType(SmallInt) }

// IsString returns true if slice is a String object
func (s Slice) IsString() bool { return s.IsType(String) }

// IsBinary returns true if slice is a Binary object
func (s Slice) IsBinary() bool { return s.IsType(Binary) }

// IsBCD returns true if slice is a BCD
func (s Slice) IsBCD() bool { return s.IsType(BCD) }

// IsCustom returns true if slice is a Custom type
func (s Slice) IsCustom() bool { return s.IsType(Custom) }

// IsInteger returns true if a slice is any decimal number type
func (s Slice) IsInteger() bool { return s.IsInt() || s.IsUInt() || s.IsSmallInt() }

// IsNumber returns true if slice is any Number-type object
func (s Slice) IsNumber() bool { return s.IsInteger() || s.IsDouble() }

// IsSorted returns true if slice is an object with table offsets, sorted by attribute name
func (s Slice) IsSorted() bool {
	h := s[0]
	return (h >= 0x0b && h <= 0x0e)
}
