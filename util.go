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

// VELOCYPACK_ASSERT panics if v is false.
func VELOCYPACK_ASSERT(v bool) {
	if !v {
		panic("VELOCYPACK_ASSERT failed")
	}
}

// read an unsigned little endian integer value of the
// specified length, starting at the specified byte offset
func readIntegerFixed(start []byte, length uint) uint64 {
	return readIntegerNonEmpty(start, length)
}

// read an unsigned little endian integer value of the
// specified length, starting at the specified byte offset
func readIntegerNonEmpty(s []byte, length uint) uint64 {
	x := uint(0)
	v := uint64(0)
	for i := uint(0); i < length; i++ {
		v += uint64(s[i]) << x
		x += 8
	}
	return v
}

func toInt64(v uint64) int64 {
	shift2 := uint64(1) << 63
	shift := int64(shift2 - 1)
	if v >= shift2 {
		return (int64(v-shift2) - shift) - 1
	} else {
		return int64(v)
	}
}

// read a variable length integer in unsigned LEB128 format
func readVariableValueLength(source []byte, reverse bool) ValueLength {
	length := ValueLength(0)
	p := uint(0)
	idx := 0
	for {
		v := ValueLength(source[idx])
		length += (v & 0x7f) << p
		p += 7
		if reverse {
			idx--
		} else {
			idx++
		}
		if v&0x80 == 0 {
			break
		}
	}
	return length
}

// optionalBool returns the first arg element if available, otherwise returns defaultValue.
func optionalBool(arg []bool, defaultValue bool) bool {
	if len(arg) == 0 {
		return defaultValue
	}
	return arg[0]
}
