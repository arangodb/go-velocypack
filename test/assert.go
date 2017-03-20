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
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func ASSERT_EQ(a, b interface{}, t *testing.T) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Expected %v, %v to be equal\nat %s", a, b, callerInfo(2))
	}
}

func ASSERT_TRUE(a bool, t *testing.T) {
	if !a {
		t.Errorf("Expected true\nat %s", callerInfo(2))
	}
}

func callerInfo(depth int) string {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "?"
	}
	return fmt.Sprintf("%s (%d)", file, line)
}
