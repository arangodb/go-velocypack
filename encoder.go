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

// This code is heavily inspired by the Go sources.
// See https://golang.org/src/encoding/json/

package velocypack

import (
	"encoding"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"sync"
)

// An Encoder encodes Go structures into velocypack values written to an output stream.
type Encoder struct {
	b Builder
	w io.Writer
}

// Marshaler is implemented by types that can convert themselves into Velocypack.
type Marshaler interface {
	MarshalVPack() ([]byte, error)
}

// NewEncoder creates a new Encoder that writes output to the given writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// Marshal writes the Velocypack encoding of v to a buffer and returns that buffer.
func Marshal(v interface{}) (result []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if s, ok := r.(string); ok {
				panic(s)
			}
			err = r.(error)
		}
	}()
	var b Builder
	reflectValue(&b, reflect.ValueOf(v))
	return b.Bytes()
}

// Encode writes the Velocypack encoding of v to the stream.
func (e *Encoder) Encode(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if s, ok := r.(string); ok {
				panic(s)
			}
			err = r.(error)
		}
	}()
	e.b.Clear()
	reflectValue(&e.b, reflect.ValueOf(v))
	if _, err := e.b.WriteTo(e.w); err != nil {
		return WithStack(err)
	}
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func reflectValue(b *Builder, v reflect.Value) {
	valueEncoder(v)(b, v)
}

type encoderFunc func(b *Builder, v reflect.Value)

var encoderCache struct {
	sync.RWMutex
	m map[reflect.Type]encoderFunc
}

func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}
	return typeEncoder(v.Type())
}

var (
	marshalerType     = reflect.TypeOf(new(Marshaler)).Elem()
	textMarshalerType = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
)

func typeEncoder(t reflect.Type) encoderFunc {
	encoderCache.RLock()
	f := encoderCache.m[t]
	encoderCache.RUnlock()
	if f != nil {
		return f
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	encoderCache.Lock()
	if encoderCache.m == nil {
		encoderCache.m = make(map[reflect.Type]encoderFunc)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	encoderCache.m[t] = func(b *Builder, v reflect.Value) {
		wg.Wait()
		f(b, v)
	}
	encoderCache.Unlock()

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = newTypeEncoder(t, true)
	wg.Done()
	encoderCache.Lock()
	encoderCache.m[t] = f
	encoderCache.Unlock()
	return f
}

// newTypeEncoder constructs an encoderFunc for a type.
// The returned encoder only checks CanAddr when allowAddr is true.
func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(marshalerType) {
			return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
		}
	}

	if t.Implements(textMarshalerType) {
		return textMarshalerEncoder
	}
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(textMarshalerType) {
			return newCondAddrEncoder(addrTextMarshalerEncoder, newTypeEncoder(t, false))
		}
	}

	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32, reflect.Float64:
		return doubleEncoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Ptr:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}

func invalidValueEncoder(b *Builder, v reflect.Value) {
	b.addNull()
}

func marshalerEncoder(b *Builder, v reflect.Value) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		b.addNull()
		return
	}
	m, ok := v.Interface().(Marshaler)
	if !ok {
		b.addNull()
		return
	}
	vpack, err := m.MarshalVPack()
	if err == nil {
		b.buf.Write(vpack)
	}
	if err != nil {
		panic(&MarshalerError{v.Type(), err})
	}
}

func addrMarshalerEncoder(b *Builder, v reflect.Value) {
	va := v.Addr()
	if va.IsNil() {
		b.addNull()
		return
	}
	m := va.Interface().(Marshaler)
	vpack, err := m.MarshalVPack()
	if err == nil {
		// copy JSON into buffer, checking validity.
		b.buf.Write(vpack)
	}
	if err != nil {
		panic(&MarshalerError{Type: v.Type(), Err: err})
	}
}

func textMarshalerEncoder(b *Builder, v reflect.Value) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		b.addNull()
		return
	}
	m := v.Interface().(encoding.TextMarshaler)
	text, err := m.MarshalText()
	if err != nil {
		panic(&MarshalerError{v.Type(), err})
	}
	b.addString(string(text))
}

func addrTextMarshalerEncoder(b *Builder, v reflect.Value) {
	va := v.Addr()
	if va.IsNil() {
		b.addNull()
		return
	}
	m := va.Interface().(encoding.TextMarshaler)
	text, err := m.MarshalText()
	if err != nil {
		panic(&MarshalerError{v.Type(), err})
	}
	b.addString(string(text))
}

func boolEncoder(b *Builder, v reflect.Value) {
	b.addBool(v.Bool())
}

func intEncoder(b *Builder, v reflect.Value) {
	b.addInt(v.Int())
}

func uintEncoder(b *Builder, v reflect.Value) {
	b.addUInt(v.Uint())
}

func doubleEncoder(b *Builder, v reflect.Value) {
	b.addDouble(v.Float())
}

func stringEncoder(b *Builder, v reflect.Value) {
	b.addString(v.String())
}

func interfaceEncoder(b *Builder, v reflect.Value) {
	if v.IsNil() {
		b.addNull()
		return
	}
	vElem := v.Elem()
	valueEncoder(vElem)(b, vElem)
}

func unsupportedTypeEncoder(b *Builder, v reflect.Value) {
	panic(&UnsupportedTypeError{v.Type()})
}

type structEncoder struct {
	fields    []field
	fieldEncs []encoderFunc
}

func (se *structEncoder) encode(b *Builder, v reflect.Value) {
	b.MustOpenObject()
	for i, f := range se.fields {
		fv := fieldByIndex(v, f.index)
		if !fv.IsValid() || f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		// Key
		b.addString(f.name)
		// Value
		se.fieldEncs[i](b, fv)
	}
	b.MustClose()
}

func newStructEncoder(t reflect.Type) encoderFunc {
	fields := cachedTypeFields(t)
	se := &structEncoder{
		fields:    fields,
		fieldEncs: make([]encoderFunc, len(fields)),
	}
	for i, f := range fields {
		se.fieldEncs[i] = typeEncoder(typeByIndex(t, f.index))
	}
	return se.encode
}

type mapEncoder struct {
	elemEnc encoderFunc
}

func (me *mapEncoder) encode(b *Builder, v reflect.Value) {
	if v.IsNil() {
		b.addNull()
	}
	b.MustOpenObject()

	// Extract and sort the keys.
	keys := v.MapKeys()
	sv := make([]reflectWithString, len(keys))
	for i, v := range keys {
		sv[i].v = v
		if err := sv[i].resolve(); err != nil {
			panic(&MarshalerError{v.Type(), err})
		}
	}
	sort.Slice(sv, func(i, j int) bool { return sv[i].s < sv[j].s })

	for _, kv := range sv {
		// Key
		b.addString(kv.s)
		// Value
		me.elemEnc(b, v.MapIndex(kv.v))
	}
	b.MustClose()
}

func newMapEncoder(t reflect.Type) encoderFunc {
	switch t.Key().Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		if !t.Key().Implements(textMarshalerType) {
			return unsupportedTypeEncoder
		}
	}
	me := &mapEncoder{typeEncoder(t.Elem())}
	return me.encode
}

func encodeByteSlice(b *Builder, v reflect.Value) {
	if v.IsNil() {
		b.addNull()
		return
	}
	b.addBinary(v.Bytes())
}

// sliceEncoder just wraps an arrayEncoder, checking to make sure the value isn't nil.
type sliceEncoder struct {
	arrayEnc encoderFunc
}

func (se *sliceEncoder) encode(b *Builder, v reflect.Value) {
	if v.IsNil() {
		b.addNull()
		return
	}
	se.arrayEnc(b, v)
}

func newSliceEncoder(t reflect.Type) encoderFunc {
	// Byte slices get special treatment; arrays don't.
	if t.Elem().Kind() == reflect.Uint8 {
		p := reflect.PtrTo(t.Elem())
		if !p.Implements(marshalerType) && !p.Implements(textMarshalerType) {
			return encodeByteSlice
		}
	}
	enc := &sliceEncoder{newArrayEncoder(t)}
	return enc.encode
}

type arrayEncoder struct {
	elemEnc encoderFunc
}

func (ae *arrayEncoder) encode(b *Builder, v reflect.Value) {
	b.MustOpenArray()
	n := v.Len()
	for i := 0; i < n; i++ {
		ae.elemEnc(b, v.Index(i))
	}
	b.MustClose()
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	enc := &arrayEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

type ptrEncoder struct {
	elemEnc encoderFunc
}

func (pe *ptrEncoder) encode(b *Builder, v reflect.Value) {
	if v.IsNil() {
		b.addNull()
		return
	}
	pe.elemEnc(b, v.Elem())
}

func newPtrEncoder(t reflect.Type) encoderFunc {
	enc := &ptrEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

type condAddrEncoder struct {
	canAddrEnc, elseEnc encoderFunc
}

func (ce *condAddrEncoder) encode(b *Builder, v reflect.Value) {
	if v.CanAddr() {
		ce.canAddrEnc(b, v)
	} else {
		ce.elseEnc(b, v)
	}
}

// newCondAddrEncoder returns an encoder that checks whether its value
// CanAddr and delegates to canAddrEnc if so, else to elseEnc.
func newCondAddrEncoder(canAddrEnc, elseEnc encoderFunc) encoderFunc {
	enc := &condAddrEncoder{canAddrEnc: canAddrEnc, elseEnc: elseEnc}
	return enc.encode
}

type reflectWithString struct {
	v reflect.Value
	s string
}

func (w *reflectWithString) resolve() error {
	if w.v.Kind() == reflect.String {
		w.s = w.v.String()
		return nil
	}
	if tm, ok := w.v.Interface().(encoding.TextMarshaler); ok {
		buf, err := tm.MarshalText()
		w.s = string(buf)
		return err
	}
	switch w.v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w.s = strconv.FormatInt(w.v.Int(), 10)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		w.s = strconv.FormatUint(w.v.Uint(), 10)
		return nil
	}
	panic("unexpected map key type")
}
