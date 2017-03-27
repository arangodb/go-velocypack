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

import "reflect"

// InvalidTypeError is returned when a Slice getter is called on a slice of a different type.
type InvalidTypeError struct {
	Message string
}

// Error implements the error interface for InvalidTypeError.
func (e InvalidTypeError) Error() string {
	return e.Message
}

// IsInvalidType returns true if the given error is an InvalidTypeError.
func IsInvalidType(err error) bool {
	_, ok := Cause(err).(InvalidTypeError)
	return ok
}

// NumberOutOfRangeError indicates an out of range error.
type NumberOutOfRangeError struct {
}

// Error implements the error interface for NumberOutOfRangeError.
func (e NumberOutOfRangeError) Error() string {
	return "number out of range"
}

// IsNumberOutOfRange returns true if the given error is an NumberOutOfRangeError.
func IsNumberOutOfRange(err error) bool {
	_, ok := Cause(err).(NumberOutOfRangeError)
	return ok
}

// IndexOutOfBoundsError indicates an index outside of array/object bounds.
type IndexOutOfBoundsError struct{}

// Error implements the error interface for IndexOutOfBoundsError.
func (e IndexOutOfBoundsError) Error() string {
	return "index out of range"
}

// IsIndexOutOfBounds returns true if the given error is an IndexOutOfBoundsError.
func IsIndexOutOfBounds(err error) bool {
	_, ok := Cause(err).(IndexOutOfBoundsError)
	return ok
}

// NeedAttributeTranslatorError indicates a lack of object key translator (smallint|uint -> string).
type NeedAttributeTranslatorError struct{}

// Error implements the error interface for NeedAttributeTranslatorError.
func (e NeedAttributeTranslatorError) Error() string {
	return "need attribute translator"
}

// IsNeedAttributeTranslator returns true if the given error is an NeedAttributeTranslatorError.
func IsNeedAttributeTranslator(err error) bool {
	_, ok := Cause(err).(NeedAttributeTranslatorError)
	return ok
}

// InternalError indicates an error that the client cannot prevent.
type InternalError struct {
}

// Error implements the error interface for InternalError.
func (e InternalError) Error() string {
	return "internal"
}

// IsInternal returns true if the given error is an InternalError.
func IsInternal(err error) bool {
	_, ok := Cause(err).(InternalError)
	return ok
}

// BuilderNeedOpenArrayError indicates an (invalid) attempt to open an array/object when that is not allowed.
type BuilderNeedOpenArrayError struct{}

// Error implements the error interface for BuilderNeedOpenArrayError.
func (e BuilderNeedOpenArrayError) Error() string {
	return "builder need open array"
}

// IsBuilderNeedOpenArray returns true if the given error is an BuilderNeedOpenArrayError.
func IsBuilderNeedOpenArray(err error) bool {
	_, ok := Cause(err).(BuilderNeedOpenArrayError)
	return ok
}

// BuilderNeedOpenObjectError indicates an (invalid) attempt to open an array/object when that is not allowed.
type BuilderNeedOpenObjectError struct{}

// Error implements the error interface for BuilderNeedOpenObjectError.
func (e BuilderNeedOpenObjectError) Error() string {
	return "builder need open object"
}

// IsBuilderNeedOpenObject returns true if the given error is an BuilderNeedOpenObjectError.
func IsBuilderNeedOpenObject(err error) bool {
	_, ok := Cause(err).(BuilderNeedOpenObjectError)
	return ok
}

// BuilderNeedOpenCompoundError indicates an (invalid) attempt to close an array/object that is already closed.
type BuilderNeedOpenCompoundError struct{}

// Error implements the error interface for BuilderNeedOpenCompoundError.
func (e BuilderNeedOpenCompoundError) Error() string {
	return "builder need open array"
}

// IsBuilderNeedOpenCompound returns true if the given error is an BuilderNeedOpenCompoundError.
func IsBuilderNeedOpenCompound(err error) bool {
	_, ok := Cause(err).(BuilderNeedOpenCompoundError)
	return ok
}

type DuplicateAttributeNameError struct{}

// Error implements the error interface for DuplicateAttributeNameError.
func (e DuplicateAttributeNameError) Error() string {
	return "duplicate key name"
}

// IsDuplicateAttributeName returns true if the given error is an DuplicateAttributeNameError.
func IsDuplicateAttributeName(err error) bool {
	_, ok := Cause(err).(DuplicateAttributeNameError)
	return ok
}

// BuilderNotSealedError is returned when a call is made to Builder.Bytes without being closed.
type BuilderNotSealedError struct{}

// Error implements the error interface for BuilderNotSealedError.
func (e BuilderNotSealedError) Error() string {
	return "builder not sealed"
}

// IsBuilderNotSealed returns true if the given error is an BuilderNotSealedError.
func IsBuilderNotSealed(err error) bool {
	_, ok := Cause(err).(BuilderNotSealedError)
	return ok
}

// BuilderKeyAlreadyWrittenError is returned when a call is made to Builder.Bytes without being closed.
type BuilderKeyAlreadyWrittenError struct{}

// Error implements the error interface for BuilderKeyAlreadyWrittenError.
func (e BuilderKeyAlreadyWrittenError) Error() string {
	return "builder key already written"
}

// IsBuilderKeyAlreadyWritten returns true if the given error is an BuilderKeyAlreadyWrittenError.
func IsBuilderKeyAlreadyWritten(err error) bool {
	_, ok := Cause(err).(BuilderKeyAlreadyWrittenError)
	return ok
}

// BuilderUnexpectedTypeError is returned when a Builder function received an invalid type.
type BuilderUnexpectedTypeError struct {
	Message string
}

// Error implements the error interface for BuilderUnexpectedTypeError.
func (e BuilderUnexpectedTypeError) Error() string {
	return e.Message
}

// IsBuilderUnexpectedType returns true if the given error is an BuilderUnexpectedTypeError.
func IsBuilderUnexpectedType(err error) bool {
	_, ok := Cause(err).(BuilderUnexpectedTypeError)
	return ok
}

// BuilderKeyMustBeStringError is returned when a key is not of type string.
type BuilderKeyMustBeStringError struct{}

// Error implements the error interface for BuilderKeyMustBeStringError.
func (e BuilderKeyMustBeStringError) Error() string {
	return "builder key must be string"
}

// IsBuilderKeyMustBeString returns true if the given error is an BuilderKeyMustBeStringError.
func IsBuilderKeyMustBeString(err error) bool {
	_, ok := Cause(err).(BuilderKeyMustBeStringError)
	return ok
}

// BuilderNeedSubValueError is returned when a RemoveLast is called without any value in an object/array.
type BuilderNeedSubValueError struct{}

// Error implements the error interface for BuilderNeedSubValueError.
func (e BuilderNeedSubValueError) Error() string {
	return "builder need sub value"
}

// IsBuilderNeedSubValue returns true if the given error is an BuilderNeedSubValueError.
func IsBuilderNeedSubValue(err error) bool {
	_, ok := Cause(err).(BuilderNeedSubValueError)
	return ok
}

// InvalidUtf8SequenceError indicates an invalid UTF8 (string) sequence.
type InvalidUtf8SequenceError struct{}

// Error implements the error interface for InvalidUtf8SequenceError.
func (e InvalidUtf8SequenceError) Error() string {
	return "invalid utf8 sequence"
}

// IsInvalidUtf8Sequence returns true if the given error is an InvalidUtf8SequenceError.
func IsInvalidUtf8Sequence(err error) bool {
	_, ok := Cause(err).(InvalidUtf8SequenceError)
	return ok
}

// MarshalerError is returned when a custom VPack Marshaler returns an error.
type MarshalerError struct {
	Type reflect.Type
	Err  error
}

// Error implements the error interface for MarshalerError.
func (e MarshalerError) Error() string {
	return "error calling MarshalVPack for type " + e.Type.String() + ": " + e.Err.Error()
}

// IsMarshaler returns true if the given error is an MarshalerError.
func IsMarshaler(err error) bool {
	_, ok := Cause(err).(MarshalerError)
	return ok
}

// UnsupportedTypeError is returned when a type is marshaled that cannot be marshaled.
type UnsupportedTypeError struct {
	Type reflect.Type
}

// Error implements the error interface for UnsupportedTypeError.
func (e UnsupportedTypeError) Error() string {
	return "unsupported type " + e.Type.String()
}

// IsUnsupportedType returns true if the given error is an UnsupportedTypeError.
func IsUnsupportedType(err error) bool {
	_, ok := Cause(err).(UnsupportedTypeError)
	return ok
}

// NoJSONEquivalentError is returned when a Velocypack type cannot be converted to JSON.
type NoJSONEquivalentError struct{}

// Error implements the error interface for NoJSONEquivalentError.
func (e NoJSONEquivalentError) Error() string {
	return "no JSON equivalent"
}

// IsNoJSONEquivalent returns true if the given error is an NoJSONEquivalentError.
func IsNoJSONEquivalent(err error) bool {
	_, ok := Cause(err).(NoJSONEquivalentError)
	return ok
}

var (
	// WithStack is called on every return of an error to add stacktrace information to the error.
	// When setting this function, also set the Cause function.
	// The interface of this function is compatible with functions in github.com/pkg/errors.
	// WithStack(nil) must return nil.
	WithStack = func(err error) error { return err }
	// Cause is used to get the root cause of the given error.
	// The interface of this function is compatible with functions in github.com/pkg/errors.
	// Cause(nil) must return nil.
	Cause = func(err error) error { return err }
)
