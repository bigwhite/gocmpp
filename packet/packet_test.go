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

package cmpppacket_test

import (
	"testing"

	"github.com/bigwhite/gocmpp/packet"
)

func TestTypeString(t *testing.T) {
	a, b, c := cmpppacket.Ver30, cmpppacket.Ver21, cmpppacket.Ver20
	if a != 0x30 || b != 0x21 || c != 0x20 {
		t.Fatal("The value of var of Type is incorrect")
	}

	if a.String() != "cmpp30" || b.String() != "cmpp21" || c.String() != "cmpp20" {
		t.Fatal("The string presentation of var of Type is incorrect")
	}
}

func TestCommandIdString(t *testing.T) {
	id1, id2 := cmpppacket.CMPP_CONNECT, cmpppacket.CMPP_CONNECT_RESP

	if id1 != 0x00000001 {
		t.Fatalf("The value of CMPP_CONNECT is %d, not equal to 0x00000001\n", id1)
	}

	if id1.String() != "CMPP_CONNECT" {
		t.Fatalf("The string presentation of command id - CMPP_CONNECT is %s, not equal to %s\n",
			id1.String(),
			"CMPP_CONNECT")
	}

	if id2 != 0x80000001 {
		t.Fatalf("The value of CMPP_CONNECT is %s, not equal to 0x00000001\n", id2)
	}
	if id2.String() != "CMPP_CONNECT_RESP" {
		t.Fatalf("The string presentation of command id - CMPP_CONNECT_RESP is %s, not equal to %s\n",
			id2.String(),
			"CMPP_CONNECT_RESP")
	}

	id3, id4 := cmpppacket.CMPP_ACTIVE_TEST, cmpppacket.CMPP_ACTIVE_TEST_RESP

	if id3 != 0x00000008 {
		t.Fatalf("The value of CMPP_ACTIVE_TEST is %d, not equal to 0x00000008\n", id3)
	}

	if id3.String() != "CMPP_ACTIVE_TEST" {
		t.Fatalf("The string presentation of command id - CMPP_ACTIVE_TEST is %s, not equal to %s\n",
			id3.String(),
			"CMPP_ACTIVE_TEST")
	}

	if id4 != 0x80000008 {
		t.Fatalf("The value of CMPP_ACTIVE_TEST_RESP is %d, not equal to 0x80000008\n", id4)
	}

	if id4.String() != "CMPP_ACTIVE_TEST_RESP" {
		t.Fatalf("The string presentation of command id - CMPP_ACTIVE_TEST_RESP is %s, not equal to %s\n",
			id4.String(),
			"CMPP_ACTIVE_TEST_RESP")
	}

	id5, id6 := cmpppacket.CMPP_GET_MO_ROUTE, cmpppacket.CMPP_GET_MO_ROUTE_RESP
	if id5 != 0x00000017 {
		t.Fatalf("The value of CMPP_GET_MO_ROUTE is %d, not equal to 0x00000017\n", id5)
	}

	if id5.String() != "CMPP_GET_MO_ROUTE" {
		t.Fatalf("The string presentation of command id - CMPP_GET_MO_ROUTE is %s, not equal to %s\n",
			id5.String(),
			"CMPP_GET_MO_ROUTE")
	}

	if id6 != 0x80000017 {
		t.Fatalf("The value of CMPP_GET_MO_ROUTE_RESP is %d, not equal to 0x80000017\n", id6)
	}

	if id6.String() != "CMPP_GET_MO_ROUTE_RESP" {
		t.Fatalf("The string presentation of command id - CMPP_GET_MO_ROUTE_RESP is %s, not equal to %s\n",
			id6.String(),
			"CMPP_GET_MO_ROUTE_RESP")
	}
}
