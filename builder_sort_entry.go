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

import "bytes"

type sortEntry struct {
	Offset ValueLength
	Name   []byte
}

type sortEntries []sortEntry

// Len is the number of elements in the collection.
func (l sortEntries) Len() int { return len(l) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (l sortEntries) Less(i, j int) bool { return bytes.Compare(l[i].Name, l[j].Name) < 0 }

// Swap swaps the elements with indexes i and j.
func (l sortEntries) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
