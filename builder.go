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

import (
	"encoding/binary"
	"math"
)

// Builder is used to build VPack structures.
type Builder struct {
	buf        builderBuffer
	stack      builderStack
	keyWritten bool
}

// Clear and start from scratch:
func (b *Builder) Clear() {
	b.buf = nil
	b.stack = nil
	b.keyWritten = false
}

// OpenObject starts a new object.
// This must be closed using Close.
func (b *Builder) OpenObject(unindexed ...bool) error {
	return nil
}

// OpenArray starts a new array.
// This must be closed using Close.
func (b *Builder) OpenArray(unindexed ...bool) error {
	return nil
}

// Close ends an open object or array.
func (b *Builder) Close() error {
	return nil
}

// IsClosed returns true if there are no more open objects or arrays.
func (b *Builder) IsClosed() bool {
	return b.stack.IsEmpty()
}

// AddNull adds a null value to the buffer.
func (b *Builder) AddNull() error {
	b.buf.WriteByte(0x18)
	return nil
}

// AddFalse adds a bool false value to the buffer.
func (b *Builder) AddFalse() error {
	b.buf.WriteByte(0x19)
	return nil
}

// AddTrue adds a bool true value to the buffer.
func (b *Builder) AddTrue() error {
	b.buf.WriteByte(0x1a)
	return nil
}

// AddDouble adds a double value to the buffer.
func (b *Builder) AddDouble(v float64) error {
	bits := math.Float64bits(v)
	binary.LittleEndian.PutUint64(b.buf.Grow(8), bits)
	return nil
}

// AddInt adds an int value to the buffer.
func (b *Builder) AddInt(v int64) error {
	if v >= 0 && v <= 9 {
		b.buf.WriteByte(0x30 + byte(v))
	} else if v < 0 && v >= -6 {
		b.buf.WriteByte(byte(0x40 + int(v)))
	} else {
		b.appendInt(v, 0x1f)
	}
	return nil
}

// AddUInt adds an uint value to the buffer.
func (b *Builder) AddUInt(v uint64) error {
	if v <= 9 {
		b.buf.WriteByte(0x30 + byte(v))
	} else {
		b.appendUInt(v, 0x27)
	}
	return nil
}

// AddUTCDate adds an UTC date value to the buffer.
func (b *Builder) AddUTCDate(v int64) error {
	x := toUInt64(v)
	b.buf.ReserveSpace(9)
	b.buf.WriteByte(0x1c)
	b.appendLength(ValueLength(x), 8)
	return nil
}

// AddString adds a string value to the buffer.
func (b *Builder) AddString(v string) error {
	raw := []byte(v)
	strLen := uint(len(raw))
	if strLen > 126 {
		// long string
		b.buf.ReserveSpace(1 + 8 + strLen)
		b.buf.WriteByte(0xbf)
		b.appendLength(ValueLength(strLen), 8) // string length
		b.buf.Write(raw)                       // string data
	} else {
		b.buf.ReserveSpace(1 + strLen)
		b.buf.WriteByte(byte(0x40 + strLen)) // short string (with length)
		b.buf.Write(raw)                     // string data
	}
	return nil
}

// AddObjectValue adds a key+value to an open object.
func (b *Builder) AddObjectValue(key string, v Value) error {
	return nil
}

// AddArrayValue adds a value to an open array.
func (b *Builder) AddArrayValue(v Value) error {
	return nil
}

// returns number of bytes required to store the value in 2s-complement
func intLength(value int64) uint {
	if value >= -0x80 && value <= 0x7f {
		// shortcut for the common case
		return 1
	}
	var x uint64
	if value >= 0 {
		x = uint64(value)
	} else {
		x = uint64(-(value + 1))
	}
	xSize := uint(0)
	for {
		xSize++
		x >>= 8
		if x < 0x80 {
			return xSize + 1
		}
	}
}

func (b *Builder) appendInt(v int64, base uint) {
	vSize := intLength(v)
	var x uint64
	if vSize == 8 {
		x = toUInt64(v)
	} else {
		shift := int64(1) << (vSize*8 - 1) // will never overflow!
		if v >= 0 {
			x = uint64(v)
		} else {
			x = uint64(v+shift) + uint64(shift)
		}
		//      x = v >= 0 ? static_cast<uint64_t>(v)
		//                 : static_cast<uint64_t>(v + shift) + shift;
	}
	b.buf.ReserveSpace(1 + vSize)
	b.buf.WriteByte(byte(base + vSize))
	for ; vSize > 0; vSize-- {
		b.buf.WriteByte(byte(x & 0xff))
		x >>= 8
	}
}

func (b *Builder) appendUInt(v uint64, base uint) {
	b.buf.ReserveSpace(9)
	save := b.buf.Len()
	b.buf.WriteByte(0) // Will be overwritten at end of function.
	vSize := uint(0)
	for {
		vSize++
		b.buf.WriteByte(byte(v & 0xff))
		v >>= 8
		if v == 0 {
			break
		}
	}
	b.buf[save] = byte(base + vSize)
}

func (b *Builder) appendLength(v ValueLength, n uint) {
	b.buf.ReserveSpace(n)
	for i := uint(0); i < n; i++ {
		b.buf.WriteByte(byte(v & 0xff))
		v >>= 8
	}
}

// openCompoundValue opens an array/object, checking the context.
func (b *Builder) openCompoundValue(vType byte) error {
	if !b.stack.IsEmpty() {
		tos := b.stack.Tos()
		buf := b.buf
		if !b.keyWritten {
			if buf[tos] != 0x06 && buf[tos] != 0x13 {
				return WithStack(BuilderNeedOpenArrayError{})
			}
		} else {
			b.keyWritten = false
		}
	}
	return nil
}

// addCompoundValue adds the start of a component value to the stream & stack.
func (b *Builder) addCompoundValue(vType byte) error {
	pos := b.buf.Len()
	b.stack.Push(pos)
	b.buf.WriteByte(vType)
	b.buf.WriteBytes(0, 8) // Will be filled later with bytelength and nr subs
	return nil
}
