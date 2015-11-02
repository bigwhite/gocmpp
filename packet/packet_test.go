// Copyright 2015 Tony Bai.
//
// Licensed under the Apache License, Vsion 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package cmpppacket

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestTypeString(t *testing.T) {
	a, b, c := V30, V21, V20
	if a != 0x30 || b != 0x21 || c != 0x20 {
		t.Fatal("The value of var of Type is incorrect")
	}

	if a.String() != "cmpp30" || b.String() != "cmpp21" || c.String() != "cmpp20" {
		t.Fatal("The string presentation of var of Type is incorrect")
	}
}

func TestCommandIdString(t *testing.T) {
	id1, id2 := CMPP_CONNECT, CMPP_CONNECT_RESP

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

	id3, id4 := CMPP_ACTIVE_TEST, CMPP_ACTIVE_TEST_RESP

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

	id5, id6 := CMPP_GET_MO_ROUTE, CMPP_GET_MO_ROUTE_RESP
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

func TestOpError(t *testing.T) {
	op := "do foo things"
	var e error = errors.New("error example for test OpError")
	err := NewOpError(e, op)

	if err.Cause() != e {
		t.Fatalf("OpError's cause: actual [%#v], wanted[%#v]\n", err.Error(), e)
	}

	if err.Op() != op {
		t.Fatalf("OpError's op : actual [%#v], wanted[%#v]\n", err.Op(), op)
	}
}

func TestPacketWriter(t *testing.T) {
	w1 := newPacketWriter(11)

	w1.WriteString("hello")
	w1.WriteString(" golang")

	s, e := w1.Bytes()
	if e != nil {
		t.Fatalf("packetWriter's err : actual [%#v], wanted[nil]\n", e)
	}

	if string(s) != "hello golang" {
		t.Fatalf("packetWriter's err : actual [%s], wanted[%s]\n", string(s), "hello golang")
	}

	var i uint16 = 0x1234
	w2 := newPacketWriter(10)
	w2.WriteInt(binary.BigEndian, i)

	s, e = w2.Bytes()
	if e != nil {
		t.Fatalf("packetWriter's err : actual [%#v], wanted[nil]\n", e)
	}

	if s[0] != 0x12 || s[1] != 0x34 {
		t.Fatalf("packetWriter's err : actual [%#v], wanted[%#v]\n", s, []byte{0x12, 0x34})
	}
}

func TestPacketReader(t *testing.T) {
	s1 := []byte{'h', 'e', 'l', 'l', 'o'}
	r1 := newPacketReader(s1)

	d1 := make([]byte, 3)
	r1.ReadBytes(d1)
	if r1.Error() != nil {
		t.Fatalf("packetReader's err : actual [%#v], wanted[nil]\n", r1.Error())
	}

	if d1[0] != 'h' || d1[1] != 'e' || d1[2] != 'l' {
		t.Fatalf("packetReader 's err : actual [%#v], wanted[%#v]\n", d1, []byte{'h', 'e', 'l'})
	}

	d2 := make([]byte, 3)
	r1.ReadBytes(d2)
	if r1.Error() == nil {
		t.Fatal("packetReader's err : actual nil, wanted non-nil")
	}

	s2 := []byte{0x12, 0x34}
	var i uint16
	r2 := newPacketReader(s2)
	r2.ReadInt(binary.BigEndian, &i)
	if r2.Error() != nil {
		t.Fatalf("packetReader's err : actual [%#v], wanted[nil]\n", r1.Error())
	}
	if i != 0x1234 {
		t.Fatalf("packetReader's err : actual [%d], wanted[%d]\n", i, 0x1234)
	}
}
