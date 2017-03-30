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

import "io"

// A Decoder decodes velocypack values into Go structures.
type Decoder struct {
	r io.Reader
}

// Unmarshaler is implemented by types that can convert themselves from Velocypack.
type Unmarshaler interface {
	UnmarshalVPack([]byte) error
}

// NewDecoder creates a new Decoder that reads data from the given reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// Unmarshal reads v from the given Velocypack encoded data slice.
func Unmarshal(data []byte, v interface{}) (err error) {
	return nil
}

// Decode reads v from the decoder stream.
func (e *Decoder) Decode(v interface{}) (err error) {
	return nil
}
