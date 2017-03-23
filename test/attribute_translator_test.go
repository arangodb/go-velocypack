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

func TestAttributeTranslator(t *testing.T) {
	tx := velocypack.NewAttributeIDTranslator()
	tx.Add("foo", 1)
	tx.Add("bar", 2)
	tx.Add("baz", 3)
	tx.Add("bark", 4)
	tx.Add("mötör", 5)
	tx.Add("quetzalcoatl", 6)

	ASSERT_EQ("foo", tx.IDToString(1), t)
	ASSERT_EQ("quetzalcoatl", tx.IDToString(6), t)
	ASSERT_EQ("", tx.IDToString(9999), t)

	ASSERT_TRUE(tx.StringToID("foo").IsSmallInt(), t)
	ASSERT_TRUE(tx.StringToID("notfound").IsNone(), t)
}
