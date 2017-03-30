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

import "sync"

var AttributeTranslator AttributeIDTranslator

// AttributeIDTranslator is used to translation integer style object keys to strings.
type AttributeIDTranslator interface {
	IDToString(id uint64) string
	StringToID(key string) Slice
}

type EditableAttributeIDTranslator interface {
	AttributeIDTranslator
	Add(key string, id uint64)
}

// attributeTranslator is a simple implementation of AttributeIDTranslator
type attributeTranslator struct {
	mutex      sync.RWMutex
	idToString map[uint64]string
	stringToID map[string]Slice
}

// NewAttributeIDTranslator creates a map based implementation of an AttributeIDTranslator.
func NewAttributeIDTranslator() EditableAttributeIDTranslator {
	return &attributeTranslator{
		idToString: make(map[uint64]string),
		stringToID: make(map[string]Slice),
	}
}

func (t *attributeTranslator) Add(key string, id uint64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.idToString[id] = key
	var b Builder
	b.addUInt(id)
	s, err := b.Slice()
	if err != nil {
		panic(err)
	}
	t.stringToID[key] = s
}

func (t *attributeTranslator) IDToString(id uint64) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if s, ok := t.idToString[id]; ok {
		return s
	}
	return ""
}

func (t *attributeTranslator) StringToID(key string) Slice {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if s, ok := t.stringToID[key]; ok {
		return s
	}
	return nil
}
