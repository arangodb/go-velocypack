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

// builderStack is a stack of positions.
type builderStack []ValueLength

// Push the given value on top of the stack
func (s *builderStack) Push(v ValueLength) {
	*s = append(*s, v)
}

// Pop removes the top of the stack.
func (s *builderStack) Pop() {
	l := len(*s)
	if l > 0 {
		*s = (*s)[:l-1]
	}
}

// Tos returns the value at the top of the stack.
// Returns <value at top of stack>, <stack is empty>
func (s builderStack) Tos() (ValueLength, bool) {
	//	_s := *s
	l := len(s)
	if l > 0 {
		return (s)[l-1], false
	}
	return 0, true
}

// IsEmpty returns true if there are no values on the stack.
func (s builderStack) IsEmpty() bool {
	l := len(s)
	return l == 0
}
