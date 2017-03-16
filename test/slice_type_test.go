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

package test

import (
	"testing"

	velocypack "github.com/arangodb/go-velocypack"
)

// TestSliceTypes checks the Type function of a slice.
func TestSliceTypes(t *testing.T) {
	s := velocypack.SliceFromHex
	tests := []struct {
		Slice velocypack.Slice
		Type  velocypack.ValueType
	}{
		{s("00"), velocypack.None},
		{s("01"), velocypack.Array},
		{s("0a"), velocypack.Object},
		{s("18"), velocypack.Null},
		{s("19"), velocypack.Bool},
		{s("1a"), velocypack.Bool},
		{s("1b"), velocypack.Double},
		{s("1c"), velocypack.UTCDate},
		{s("1e"), velocypack.MinKey},
		{s("1f"), velocypack.MaxKey},
	}

	for _, test := range tests {
		vt := test.Slice.Type()
		if vt != test.Type {
			t.Errorf("Invalid type for '%s', expected '%s', got '%s'", test.Slice, test.Type, vt)
		}
	}
}
