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

var (
	// WithStack is called on every return of an error to add stacktrace information to the error.
	// When setting this function, also set the Cause function.
	// The interface of this function is compatible with functions in github.com/pkg/errors.
	WithStack = func(err error) error { return err }
	// Cause is used to get the root cause of the given error.
	// The interface of this function is compatible with functions in github.com/pkg/errors.
	Cause = func(err error) error { return err }
)
