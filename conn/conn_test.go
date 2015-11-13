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

package cmppconn_test

import (
	"testing"

	cmppconn "github.com/bigwhite/gocmpp/conn"
)

func TestTypeString(t *testing.T) {
	a, b, c := cmppconn.V30, cmppconn.V21, cmppconn.V20
	if a != 0x30 || b != 0x21 || c != 0x20 {
		t.Fatal("The value of var of Type is incorrect")
	}

	if a.String() != "cmpp30" || b.String() != "cmpp21" || c.String() != "cmpp20" {
		t.Fatal("The string presentation of var of Type is incorrect")
	}
}
