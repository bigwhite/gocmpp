// Copyright 2015 Tony Bai.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package cmpputils_test

import (
	"testing"

	"github.com/bigwhite/gocmpp/utils"
)

func TestTimeStamp2Str(t *testing.T) {
	var t1 uint32 = 1021080510
	s1 := cmpputils.TimeStamp2Str(t1)
	if s1 != "1021080510" {
		t.Errorf("The result of TimeStamp2Str is %s, not equal to expected: %s\n", s1, "1021080510")
	}

	var t2 uint32 = 121080510
	s2 := cmpputils.TimeStamp2Str(t2)
	if s2 != "0121080510" {
		t.Errorf("The result of TimeStamp2Str is %s, not equal to expected: %s\n", s2, "0121080510")
	}
}
