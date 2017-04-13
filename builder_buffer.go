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

// builderBuffer is a byte slice used for building slices.
type builderBuffer []byte

const (
	minGrowDelta = 32          // Minimum amount of extra bytes to add to a buffer when growing
	maxGrowDelta = 1024 * 1024 // Maximum amount of extra bytes to add to a buffer when growing
)

// IsEmpty returns 0 if there are no values in the buffer.
func (b *builderBuffer) IsEmpty() bool {
	l := len(*b)
	return l == 0
}

// Len returns the length of the buffer.
func (b *builderBuffer) Len() ValueLength {
	l := len(*b)
	return ValueLength(l)
}

// Bytes returns the bytes written to the buffer.
// The returned slice is only valid until the next modification.
func (b *builderBuffer) Bytes() []byte {
	return *b
}

// WriteByte appends a single byte to the buffer.
func (b *builderBuffer) WriteByte(v byte) {
	off := len(*b)
	b.grow(1)
	(*b)[off] = v
}

// WriteBytes appends a series of identical bytes to the buffer.
func (b *builderBuffer) WriteBytes(v byte, count uint) {
	off := uint(len(*b))
	b.grow(count)
	for i := uint(0); i < count; i++ {
		(*b)[off+i] = v
	}
}

// Write appends a series of bytes to the buffer.
func (b *builderBuffer) Write(v []byte) {
	l := len(v)
	if l > 0 {
		off := len(*b)
		b.grow(uint(l))
		copy((*b)[off:], v)
	}
}

// ReserveSpace ensures that at least n bytes can be added to the buffer without allocating new memory.
func (b *builderBuffer) ReserveSpace(n uint) {
	if n > 0 {
		l := len(*b)
		b.grow(n)
		*b = (*b)[:l]
	}
}

// Shrink reduces the length of the buffer by n elements (removing the last elements).
func (b *builderBuffer) Shrink(n uint) {
	if n > 0 {
		newLen := uint(len(*b)) - n
		if newLen < 0 {
			newLen = 0
		}
		*b = (*b)[0:newLen]
	}
}

// Grow adds n elements to the buffer, returning a slice where the added elements start.
func (b *builderBuffer) Grow(n uint) []byte {
	l := len(*b)
	if n > 0 {
		b.grow(n)
	}
	return (*b)[l:]
}

// grow adds n elements to the buffer.
func (b *builderBuffer) grow(n uint) {
	var newBuffer builderBuffer
	newLen := uint(len(*b)) + n
	if newLen > uint(cap(*b)) {
		extra := newLen / 4
		if extra < minGrowDelta {
			extra = minGrowDelta
		} else if extra > maxGrowDelta {
			extra = maxGrowDelta
		}
		newBuffer = make(builderBuffer, newLen, newLen+extra)
		copy(newBuffer, *b)
	} else {
		newBuffer = (*b)[0:newLen]
	}
	*b = newBuffer
}
