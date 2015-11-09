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

	cmpppacket "github.com/bigwhite/gocmpp/packet"
)

func TestCmppActiveTestReqPktPack(t *testing.T) {
	p := &cmpppacket.CmppActiveTestReqPkt{}

	data, err := p.Pack(seqId)
	if err != nil {
		t.Fatal("CmppActiveTestReqPkt pack error:", err)
	}

	if p.SeqId != seqId {
		t.Fatalf("After pack, seqId is %d, not equal to expected: %d\n", p.SeqId, seqId)
	}

	// data after pack expected:
	dataExpected := []byte{
		0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x17,
	}

	l1 := len(data)
	l2 := len(dataExpected)
	if l1 != l2 {
		t.Fatalf("After pack, data length is %d, not equal to length expected: %d\n", l1, l2)
	}

	for i := 0; i < l1; i++ {
		if data[i] != dataExpected[i] {
			t.Fatalf("After pack, data[%d] is %x, not equal to dataExpected[%d]: %x\n", i, data[i], i, dataExpected[i])
		}
	}
}

func TestCmppActiveTestReqUnpack(t *testing.T) {
	// cmpp active test request packet data:
	data := []byte{
		0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x17,
	}

	p := &cmpppacket.CmppActiveTestReqPkt{}
	err := p.Unpack(data[8:])
	if err != nil {
		t.Fatal("CmppActiveTestReqPkt unpack error:", err)
	}

	if p.SeqId != seqId {
		t.Fatalf("After unpack, seqId in packet is %x, not equal to the expected value: %x\n", p.SeqId, seqId)
	}
}

func TestCmppActiveTestRspPktPack(t *testing.T) {
	p := &cmpppacket.CmppActiveTestRspPkt{}

	data, err := p.Pack(seqId)
	if err != nil {
		t.Fatal("CmppActiveTestRspPkt pack error:", err)
	}

	if p.SeqId != seqId {
		t.Fatalf("After pack, seqId is %d, not equal to expected: %d\n", p.SeqId, seqId)
	}

	// data after pack expected:
	dataExpected := []byte{
		0x00, 0x00, 0x00, 0x0d, 0x80, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x17, 0x00,
	}

	l1 := len(data)
	l2 := len(dataExpected)
	if l1 != l2 {
		t.Fatalf("After pack, data length is %d, not equal to length expected: %d\n", l1, l2)
	}

	for i := 0; i < l1; i++ {
		if data[i] != dataExpected[i] {
			t.Fatalf("After pack, data[%d] is %x, not equal to dataExpected[%d]: %x\n", i, data[i], i, dataExpected[i])
		}
	}
}

func TestCmppActiveTestRspUnpack(t *testing.T) {
	// cmpp active test response packet data:
	data := []byte{
		0x00, 0x00, 0x00, 0x0d, 0x80, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x17, 0x00,
	}

	p := &cmpppacket.CmppActiveTestRspPkt{}
	err := p.Unpack(data[8:])
	if err != nil {
		t.Fatal("CmppActiveTestRspPkt unpack error:", err)
	}

	if p.SeqId != seqId {
		t.Fatalf("After unpack, seqId in packet is %x, not equal to the expected value: %x\n", p.SeqId, seqId)
	}
}
