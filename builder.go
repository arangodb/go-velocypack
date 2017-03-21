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
	"sort"
)

// BuilderOptions contains options that influence how Builder builds slices.
type BuilderOptions struct {
	BuildUnindexedArrays     bool
	BuildUnindexedObjects    bool
	CheckAttributeUniqueness bool
}

// Builder is used to build VPack structures.
type Builder struct {
	BuilderOptions
	buf        builderBuffer
	stack      builderStack
	index      []indexVector
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
	var vType byte
	if optionalBool(unindexed, false) {
		vType = 0x14
	} else {
		vType = 0x0b
	}
	return WithStack(b.openCompoundValue(vType))
}

// OpenArray starts a new array.
// This must be closed using Close.
func (b *Builder) OpenArray(unindexed ...bool) error {
	var vType byte
	if optionalBool(unindexed, false) {
		vType = 0x13
	} else {
		vType = 0x06
	}
	return WithStack(b.openCompoundValue(vType))
}

// Close ends an open object or array.
func (b *Builder) Close() error {
	if b.IsClosed() {
		return WithStack(BuilderNeedOpenCompoundError{})
	}
	tos := b.stack.Tos()
	head := b.buf[tos]

	VELOCYPACK_ASSERT(head == 0x06 || head == 0x0b || head == 0x13 ||
		head == 0x14)

	isArray := (head == 0x06 || head == 0x13)
	index := b.index[len(b.stack)-1]

	if index.IsEmpty() {
		b.closeEmptyArrayOrObject(tos, isArray)
		return nil
	}

	// From now on index.size() > 0
	VELOCYPACK_ASSERT(len(index) > 0)

	// check if we can use the compact Array / Object format
	if head == 0x13 || head == 0x14 ||
		(head == 0x06 && b.BuilderOptions.BuildUnindexedArrays) ||
		(head == 0x0b && (b.BuilderOptions.BuildUnindexedObjects || len(index) == 1)) {
		if b.closeCompactArrayOrObject(tos, isArray, index) {
			return nil
		}
		// This might fall through, if closeCompactArrayOrObject gave up!
	}

	if isArray {
		b.closeArray(tos, index)
		return nil
	}

	// From now on we're closing an object

	// fix head byte in case a compact Array / Object was originally requested
	b.buf[tos] = 0x0b

	// First determine byte length and its format:
	offsetSize := uint(8)
	// can be 1, 2, 4 or 8 for the byte width of the offsets,
	// the byte length and the number of subvalues:
	if b.buf.Len()-tos+ValueLength(len(index))-6 <= 0xff {
		// We have so far used _pos - tos bytes, including the reserved 8
		// bytes for byte length and number of subvalues. In the 1-byte number
		// case we would win back 6 bytes but would need one byte per subvalue
		// for the index table
		offsetSize = 1

		// Maybe we need to move down data:
		targetPos := ValueLength(3)
		if b.buf.Len() > (tos + 9) {
			_len := ValueLength(b.buf.Len() - (tos + 9))
			checkOverflow(_len)
			src := b.buf[tos+9:]
			copy(b.buf[tos+targetPos:], src[:_len])
		}
		diff := ValueLength(9 - targetPos)
		b.buf.Shrink(uint(diff))
		n := len(index)
		for i := 0; i < n; i++ {
			index[i] -= diff
		}

		// One could move down things in the offsetSize == 2 case as well,
		// since we only need 4 bytes in the beginning. However, saving these
		// 4 bytes has been sacrificed on the Altar of Performance.
	} else if b.buf.Len()-tos+2*ValueLength(len(index)) <= 0xffff {
		offsetSize = 2
	} else if b.buf.Len()-tos+4*ValueLength(len(index)) <= 0xffffffff {
		offsetSize = 4
	}

	// Now build the table:
	extraSpace := offsetSize * uint(len(index))
	if offsetSize == 8 {
		extraSpace += 8
	}
	b.buf.ReserveSpace(extraSpace)
	tableBase := b.buf.Len()
	b.buf.Grow(offsetSize * uint(len(index)))
	// Object
	if len(index) >= 2 {
		if err := b.sortObjectIndex(b.buf[tos:], index); err != nil {
			return WithStack(err)
		}
	}
	for i := uint(0); i < uint(len(index)); i++ {
		indexBase := tableBase + ValueLength(offsetSize*i)
		x := uint64(index[i])
		for j := uint(0); j < offsetSize; j++ {
			b.buf[indexBase+ValueLength(j)] = byte(x & 0xff)
			x >>= 8
		}
	}
	// Finally fix the byte width in the type byte:
	if offsetSize > 1 {
		if offsetSize == 2 {
			b.buf[tos] += 1
		} else if offsetSize == 4 {
			b.buf[tos] += 2
		} else { // offsetSize == 8
			b.buf[tos] += 3
			b.appendLength(ValueLength(len(index)), 8)
		}
	}

	// Fix the byte length in the beginning:
	x := ValueLength(b.buf.Len() - tos)
	for i := uint(1); i <= offsetSize; i++ {
		b.buf[tos+ValueLength(i)] = byte(x & 0xff)
		x >>= 8
	}

	if offsetSize < 8 {
		x := len(index)
		for i := uint(offsetSize + 1); i <= 2*offsetSize; i++ {
			b.buf[tos+ValueLength(i)] = byte(x & 0xff)
			x >>= 8
		}
	}

	// And, if desired, check attribute uniqueness:
	if b.BuilderOptions.CheckAttributeUniqueness && len(index) > 1 {
		// check uniqueness of attribute names
		if err := b.checkAttributeUniqueness(Slice(b.buf[tos:])); err != nil {
			return WithStack(err)
		}
	}

	// Now the array or object is complete, we pop a ValueLength off the _stack:
	b.stack.Pop()
	// Intentionally leave _index[depth] intact to avoid future allocs!
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
	b.addCompoundValue(vType)
	return nil
}

// addCompoundValue adds the start of a component value to the stream & stack.
func (b *Builder) addCompoundValue(vType byte) {
	pos := b.buf.Len()
	b.stack.Push(pos)
	for len(b.stack) < len(b.index) {
		b.index = append(b.index, indexVector{})
	}
	b.index[len(b.stack)-1].Clear()
	b.buf.WriteByte(vType)
	b.buf.WriteBytes(0, 8) // Will be filled later with bytelength and nr subs
}

// closeEmptyArrayOrObject closes an empty array/object, removing the pre-allocated length space.
func (b *Builder) closeEmptyArrayOrObject(tos ValueLength, isArray bool) {
	// empty Array or Object
	if isArray {
		b.buf[tos] = 0x01
	} else {
		b.buf[tos] = 0x0a
	}
	VELOCYPACK_ASSERT(b.buf.Len() == tos+9)
	b.buf.Shrink(8)
	b.stack.Pop()
}

// closeCompactArrayOrObject tries to close an array/object using compact notation.
// Returns true when a compact notation was possible, false otherwise.
func (b *Builder) closeCompactArrayOrObject(tos ValueLength, isArray bool, index indexVector) bool {
	// use compact notation
	nrItems := len(index)
	nrItemsLen := getVariableValueLength(ValueLength(nrItems))
	VELOCYPACK_ASSERT(nrItemsLen > 0)

	byteSize := b.buf.Len() - (tos + 8) + nrItemsLen
	VELOCYPACK_ASSERT(byteSize > 0)

	byteSizeLen := getVariableValueLength(byteSize)
	byteSize += byteSizeLen
	if getVariableValueLength(byteSize) != byteSizeLen {
		byteSize++
		byteSizeLen++
	}

	if byteSizeLen < 9 {
		// can only use compact notation if total byte length is at most 8 bytes long
		if isArray {
			b.buf[tos] = 0x13
		} else {
			b.buf[tos] = 0x14
		}

		valuesLen := b.buf.Len() - (tos + 9) // Amount of bytes taken up by array/object values.
		if valuesLen > 0 && byteSizeLen < 8 {
			// We have array/object values and our byteSize needs less than the pre-allocated 8 bytes.
			// So we move the array/object values back.
			checkOverflow(valuesLen)
			src := b.buf[tos+9:]
			copy(b.buf[tos+1+byteSizeLen:], src[:valuesLen])
		}
		// Shrink buffer, removing unused space allocated for byteSize.
		b.buf.Shrink(uint(8 - byteSizeLen))

		// store byte length
		VELOCYPACK_ASSERT(byteSize > 0)
		storeVariableValueLength(b.buf, tos+1, byteSize, false)

		// store nrItems
		nrItemsDst := b.buf.Grow(uint(nrItemsLen))
		storeVariableValueLength(nrItemsDst, 0, ValueLength(len(index)), false)

		b.stack.Pop()
		return true
	}
	return false
}

// checkAttributeUniqueness checks the given slice for duplicate keys.
// It returns an error when duplicate keys are found, nil otherwise.
func (b *Builder) checkAttributeUniqueness(obj Slice) error {
	VELOCYPACK_ASSERT(b.BuilderOptions.CheckAttributeUniqueness)
	n, err := obj.Length()
	if err != nil {
		return WithStack(err)
	}

	if obj.IsSorted() {
		// object attributes are sorted
		previous, err := obj.KeyAt(0)
		if err != nil {
			return WithStack(err)
		}
		p, err := previous.GetString()
		if err != nil {
			return WithStack(err)
		}

		// compare each two adjacent attribute names
		for i := ValueLength(1); i < n; i++ {
			current, err := obj.KeyAt(i)
			if err != nil {
				return WithStack(err)
			}
			// keyAt() guarantees a string as returned type
			VELOCYPACK_ASSERT(current.IsString())

			q, err := current.GetString()
			if err != nil {
				return WithStack(err)
			}

			if p == q {
				// identical key
				return WithStack(DuplicateAttributeNameError{})
			}
			// re-use already calculated values for next round
			p = q
		}
	} else {
		keys := make(map[string]struct{})

		for i := ValueLength(0); i < n; i++ {
			// note: keyAt() already translates integer attributes
			key, err := obj.KeyAt(i)
			if err != nil {
				return WithStack(err)
			}
			// keyAt() guarantees a string as returned type
			VELOCYPACK_ASSERT(key.IsString())

			k, err := key.GetString()
			if err != nil {
				return WithStack(err)
			}
			if _, found := keys[k]; found {
				return WithStack(DuplicateAttributeNameError{})
			}
			keys[k] = struct{}{}
		}
	}
	return nil
}

func findAttrName(base []byte) ([]byte, error) {
	b := base[0]
	if b >= 0x40 && b <= 0xbe {
		// short UTF-8 string
		l := b - 0x40
		return base[1 : 1+l], nil
	}
	if b == 0xbf {
		// long UTF-8 string
		l := uint(0)
		// read string length
		for i := 8; i >= 1; i-- {
			l = (l << 8) + uint(base[i])
		}
		return base[1+8 : 1+8+l], nil
	}

	// translate attribute name
	key, err := Slice(base).makeKey()
	if err != nil {
		return nil, WithStack(err)
	}
	return findAttrName(key)
}

func (b *Builder) sortObjectIndex(objBase []byte, offsets []ValueLength) error {
	var list sortEntries
	for _, off := range offsets {
		name, err := findAttrName(objBase[off:])
		if err != nil {
			return WithStack(err)
		}
		list = append(list, sortEntry{
			Offset: off,
			Name:   name,
		})
	}
	sort.Sort(list)
	for i, entry := range list {
		offsets[i] = entry.Offset
	}
	return nil
}

func (b *Builder) closeArray(tos ValueLength, index []ValueLength) {
	// fix head byte in case a compact Array was originally requested:
	b.buf[tos] = 0x06

	needIndexTable := true
	needNrSubs := true
	if len(index) == 1 {
		needIndexTable = false
		needNrSubs = false
	} else if (b.buf.Len()-tos)-index[0] == ValueLength(len(index))*(index[1]-index[0]) {
		// In this case it could be that all entries have the same length
		// and we do not need an offset table at all:
		noTable := true
		subLen := index[1] - index[0]
		if (b.buf.Len()-tos)-index[len(index)-1] != subLen {
			noTable = false
		} else {
			for i := 1; i < len(index)-1; i++ {
				if index[i+1]-index[i] != subLen {
					noTable = false
					break
				}
			}
		}
		if noTable {
			needIndexTable = false
			needNrSubs = false
		}
	}

	// First determine byte length and its format:
	var offsetSize uint
	// can be 1, 2, 4 or 8 for the byte width of the offsets,
	// the byte length and the number of subvalues:
	var indexLenIfNeeded ValueLength
	if needIndexTable {
		indexLenIfNeeded = ValueLength(len(index))
	}
	nrSubsLenIfNeeded := ValueLength(7)
	if needNrSubs {
		nrSubsLenIfNeeded = 6
	}
	if b.buf.Len()-tos+(indexLenIfNeeded)-(nrSubsLenIfNeeded) <= 0xff {
		// We have so far used _pos - tos bytes, including the reserved 8
		// bytes for byte length and number of subvalues. In the 1-byte number
		// case we would win back 6 bytes but would need one byte per subvalue
		// for the index table
		offsetSize = 1
	} else if b.buf.Len()-tos+(indexLenIfNeeded*2) <= 0xffff {
		offsetSize = 2
	} else if b.buf.Len()-tos+(indexLenIfNeeded*4) <= 0xffffffff {
		offsetSize = 4
	} else {
		offsetSize = 8
	}

	// Maybe we need to move down data:
	if offsetSize == 1 {
		targetPos := ValueLength(3)
		if !needIndexTable {
			targetPos = 2
		}
		if b.buf.Len() > (tos + 9) {
			_len := ValueLength(b.buf.Len() - (tos + 9))
			checkOverflow(_len)
			src := b.buf[tos+9:]
			copy(b.buf[tos+targetPos:], src[:_len])
		}
		diff := ValueLength(9 - targetPos)
		b.buf.Shrink(uint(diff))
		if needIndexTable {
			n := len(index)
			for i := 0; i < n; i++ {
				index[i] -= diff
			}
		} // Note: if !needIndexTable the index array is now wrong!
	}
	// One could move down things in the offsetSize == 2 case as well,
	// since we only need 4 bytes in the beginning. However, saving these
	// 4 bytes has been sacrificed on the Altar of Performance.

	// Now build the table:
	if needIndexTable {
		extraSpaceNeeded := offsetSize * uint(len(index))
		if offsetSize == 8 {
			extraSpaceNeeded += 8
		}
		b.buf.ReserveSpace(extraSpaceNeeded)
		tableBase := b.buf.Grow(offsetSize * uint(len(index)))
		for i := uint(0); i < uint(len(index)); i++ {
			x := uint64(index[i])
			for j := uint(0); j < offsetSize; j++ {
				tableBase[offsetSize*i+j] = byte(x & 0xff)
				x >>= 8
			}
		}
	} else { // no index table
		b.buf[tos] = 0x02
	}
	// Finally fix the byte width in the type byte:
	if offsetSize > 1 {
		if offsetSize == 2 {
			b.buf[tos] += 1
		} else if offsetSize == 4 {
			b.buf[tos] += 2
		} else { // offsetSize == 8
			b.buf[tos] += 3
			if needNrSubs {
				b.appendLength(ValueLength(len(index)), 8)
			}
		}
	}

	// Fix the byte length in the beginning:
	x := ValueLength(b.buf.Len() - tos)
	for i := uint(1); i <= offsetSize; i++ {
		b.buf[tos+ValueLength(i)] = byte(x & 0xff)
		x >>= 8
	}

	if offsetSize < 8 && needNrSubs {
		x = ValueLength(len(index))
		for i := offsetSize + 1; i <= 2*offsetSize; i++ {
			b.buf[tos+ValueLength(i)] = byte(x & 0xff)
			x >>= 8
		}
	}

	// Now the array or object is complete, we pop a ValueLength
	// off the _stack:
	b.stack.Pop()
	// Intentionally leave _index[depth] intact to avoid future allocs!
}
