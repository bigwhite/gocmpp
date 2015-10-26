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

func TestConnectRequestPacketPack(t *testing.T) {
	p := &cmpppacket.ConnectRequestPacket{
		SourceAddr: "900001",
		Version:    cmpppacket.Ver21,
		Timestamp:  1021080510,
		Secret:     "888888",
	}

	var seqId uint32 = 0x17
	data, err := p.Pack(seqId)
	if err != nil {
		t.Fatal("ConnectRequestPacket pack error:", err)
	}

	if p.SeqId != seqId {
		t.Fatalf("After pack, seqId is %d, not equal to expected: %d\n", p.SeqId, seqId)
	}

	// data after pack expected:
	//
	// 00000000  00 00 00 27 00 00 00 01  00 00 00 17 39 30 30 30  |...'........9000|
	// 00000010  30 31 90 d0 0c 1d 51 7a  bd 0b 4f 65 f6 bc f8 53  |01....Qz..Oe...S|
	// 00000020  5d 16 21 3c dc 73 be                              |].!<.s.|
	dataExpected := []byte{
		0x00, 0x00, 0x00, 0x27, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x17, 0x39, 0x30, 0x30, 0x30,
		0x30, 0x31, 0x90, 0xd0, 0x0c, 0x1d, 0x51, 0x7a, 0xbd, 0x0b, 0x4f, 0x65, 0xf6, 0xbc, 0xf8, 0x53,
		0x5d, 0x16, 0x21, 0x3c, 0xdc, 0x73, 0xbe,
	}

	l1 := len(data)
	l2 := len(dataExpected)
	if l1 != l2 {
		t.Fatalf("After pack, data length is %d, not equal to length expected: %d\n", l1, l2)
	}

	for i := 0; i < l1; i++ {
		if data[i] != dataExpected[i] {
			t.Fatalf("After pack, data[%d] is %x, not equal to dataExpected[%d]: %x\n", data[i], dataExpected[i])
		}
	}
}

func TestConnectRequestPacketUnpack(t *testing.T) {
	// connect request packet data:
	//
	// 00000000  00 00 00 27 00 00 00 01  00 00 00 17 39 30 30 30  |...'........9000|
	// 00000010  30 31 90 d0 0c 1d 51 7a  bd 0b 4f 65 f6 bc f8 53  |01....Qz..Oe...S|
	// 00000020  5d 16 21 3c dc 73 be                              |].!<.s.|
	data := []byte{
		0x00, 0x00, 0x00, 0x27, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x17, 0x39, 0x30, 0x30, 0x30,
		0x30, 0x31, 0x90, 0xd0, 0x0c, 0x1d, 0x51, 0x7a, 0xbd, 0x0b, 0x4f, 0x65, 0xf6, 0xbc, 0xf8, 0x53,
		0x5d, 0x16, 0x21, 0x3c, 0xdc, 0x73, 0xbe,
	}

	p := &cmpppacket.ConnectRequestPacket{}
	_ = p.Unpack(data[8:])

	if p.SeqId != 0x17 {
		t.Fatalf("After unpack, seqId in packet is %x, not equal to the expected value: 0x17\n", p.SeqId)
	}
	if p.SourceAddr != "900001" {
		t.Fatalf("After unpack, SourceAddr in packet is %s, not equal to the expected value: 900001\n", p.SourceAddr)
	}
	if p.Version != cmpppacket.Ver21 {
		t.Fatalf("After unpack, Version in packet is %x, not equal to the expected value: %x\n",
			p.Version, cmpppacket.Ver21)
	}
	if p.Timestamp != 1021080510 {
		t.Fatalf("After unpack, Timestamp in packet is %d, not equal to the expected value: 1021080510\n", p.Timestamp)
	}
}
