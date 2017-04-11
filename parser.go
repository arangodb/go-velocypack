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
	"encoding/json"
	"io"
	"math"
	"strings"
)

// Parser is used to build VPack structures from JSON.
type Parser struct {
	decoder *json.Decoder
	builder *Builder
}

// ParseJSON parses JSON from the given reader and returns the
// VPack equivalent.
func ParseJSON(r io.Reader) (Slice, error) {
	builder := &Builder{}
	p := NewParser(r, builder)
	if err := p.Parse(); err != nil {
		return nil, WithStack(err)
	}
	slice, err := builder.Slice()
	if err != nil {
		return nil, WithStack(err)
	}
	return slice, nil
}

// ParseJSONFromString parses the given JSON string and returns the
// VPack equivalent.
func ParseJSONFromString(json string) (Slice, error) {
	return ParseJSON(strings.NewReader(json))
}

// NewParser initializes a new Parser with JSON from the given reader and
// it will store the parsers output in the given builder.
func NewParser(r io.Reader, builder *Builder) *Parser {
	return &Parser{
		decoder: json.NewDecoder(r),
		builder: builder,
	}
}

// Parse JSON from the parsers reader and build VPack structures in the
// parsers builder.
func (p *Parser) Parse() error {
	for {
		t, err := p.decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return WithStack(err)
		}
		switch x := t.(type) {
		case nil:
			if err := p.builder.AddValue(NewNullValue()); err != nil {
				return WithStack(err)
			}
		case bool:
			if err := p.builder.AddValue(NewBoolValue(x)); err != nil {
				return WithStack(err)
			}
		case float64:
			if math.Trunc(x) == x {
				// It's an integer
				if x < 0 {
					if x >= math.MinInt64 {
						if err := p.builder.AddValue(NewIntValue(int64(x))); err != nil {
							return WithStack(err)
						}
					} else {
						// Does not fit in int64
						if err := p.builder.AddValue(NewDoubleValue(x)); err != nil {
							return WithStack(err)
						}
					}
				} else {
					if x <= math.MaxUint64 {
						if err := p.builder.AddValue(NewUIntValue(uint64(x))); err != nil {
							return WithStack(err)
						}
					} else {
						// Does not fit in uint64
						if err := p.builder.AddValue(NewDoubleValue(x)); err != nil {
							return WithStack(err)
						}
					}
				}
			} else {
				// Floating point
				if err := p.builder.AddValue(NewDoubleValue(x)); err != nil {
					return WithStack(err)
				}
			}
		case string:
			if err := p.builder.AddValue(NewStringValue(x)); err != nil {
				return WithStack(err)
			}
		case json.Delim:
			switch x {
			case '[':
				if err := p.builder.OpenArray(); err != nil {
					return WithStack(err)
				}
			case '{':
				if err := p.builder.OpenObject(); err != nil {
					return WithStack(err)
				}
			case ']', '}':
				if err := p.builder.Close(); err != nil {
					return WithStack(err)
				}
			}
		}
	}
	return nil
}
