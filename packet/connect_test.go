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

var (
	sourceAddr        = "900001"
	secret            = "888888"
	timestamp  uint32 = 1021080510
	version           = cmpppacket.V21
	seqId      uint32 = 0x17
)

func TestCmppConnReqPktPack(t *testing.T) {
	p := &cmpppacket.CmppConnReqPkt{
		SrcAddr:   sourceAddr,
		Version:   version,
		Secret:    secret,
		Timestamp: timestamp, //usually , we don't need to assign timestamp
	}

	data, err := p.Pack(seqId)
	if err != nil {
		t.Fatal("CmppConnReqPkt pack error:", err)
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
			t.Fatalf("After pack, data[%d] is %x, not equal to dataExpected[%d]: %x\n", i, data[i], i, dataExpected[i])
		}
	}
}

func TestCmppConnReqPktUnpack(t *testing.T) {
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

	p := &cmpppacket.CmppConnReqPkt{}
	err := p.Unpack(data[8:])
	if err != nil {
		t.Fatal("CmppConnReqPkt unpack error:", err)
	}

	if p.SeqId != seqId {
		t.Fatalf("After unpack, seqId in packet is %x, not equal to the expected value: %x\n", p.SeqId, seqId)
	}
	if p.SrcAddr != sourceAddr {
		t.Fatalf("After unpack, SrcAddr in packet is %s, not equal to the expected value: %s\n", p.SrcAddr, sourceAddr)
	}
	if p.Version != version {
		t.Fatalf("After unpack, Version in packet is %x, not equal to the expected value: %x\n",
			p.Version, version)
	}
	if p.Timestamp != timestamp {
		t.Fatalf("After unpack, Timestamp in packet is %d, not equal to the expected value: %d\n", p.Timestamp, timestamp)
	}

	authSrcExpected := []byte{
		0x90, 0xd0, 0x0c, 0x1d, 0x51, 0x7a, 0xbd, 0x0b,
		0x4f, 0x65, 0xf6, 0xbc, 0xf8, 0x53, 0x5d, 0x16,
	}

	authSrc := []byte(p.AuthSrc)
	for i := 0; i < len(authSrc); i++ {
		if authSrc[i] != authSrcExpected[i] {
			t.Fatalf("After unpack, authSrc[%d] is %x, not equal to authsrcExpected[%d]: %x\n", i, authSrc[i], i, authSrcExpected[i])
		}
	}
}

func TestCmpp2ConnRspPktPack(t *testing.T) {
	//AuthSrc: 90 d0 0c 1d 51 7a bd 0b  4f 65 f6 bc f8 53 5d 16
	authSrc := []byte{
		0x90, 0xd0, 0x0c, 0x1d, 0x51, 0x7a, 0xbd, 0x0b,
		0x4f, 0x65, 0xf6, 0xbc, 0xf8, 0x53, 0x5d, 0x16,
	}

	p := &cmpppacket.Cmpp2ConnRspPkt{
		Status:  0x0,
		Version: version,
		Secret:  secret,
		AuthSrc: string(authSrc),
	}

	data, err := p.Pack(seqId)
	if err != nil {
		t.Fatal("Cmpp2ConnRspPkt pack error:", err)
	}

	if p.SeqId != seqId {
		t.Fatalf("After pack, seqId is %d, not equal to expected: %d\n", p.SeqId, seqId)
	}

	// data after pack expected
	dataExpected := []byte{
		0x00, 0x00, 0x00, 0x1e, 0x80, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x17, 0x00,
		0x6c, 0x0b, 0x84, 0x6e, 0x25, 0xba, 0xb6, 0xda, 0xa4, 0xed, 0x1c, 0x46, 0x6e,
		0x0f, 0x4b, 0xd8, 0x21,
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

func TestCmpp2ConnRspUnpack(t *testing.T) {
	// cmpp2 connect response packet data:
	data := []byte{
		0x00, 0x00, 0x00, 0x1e, 0x80, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x17, 0x00,
		0x6c, 0x0b, 0x84, 0x6e, 0x25, 0xba, 0xb6, 0xda, 0xa4, 0xed, 0x1c, 0x46, 0x6e,
		0x0f, 0x4b, 0xd8, 0x21,
	}

	p := &cmpppacket.Cmpp2ConnRspPkt{}
	err := p.Unpack(data[8:])
	if err != nil {
		t.Fatal("Cmpp2ConnRspPkt unpack error:", err)
	}

	if p.SeqId != seqId {
		t.Fatalf("After unpack, seqId in packet is %x, not equal to the expected value: %x\n", p.SeqId, seqId)
	}
	if p.Version != version {
		t.Fatalf("After unpack, Version in packet is %x, not equal to the expected value: %x\n",
			p.Version, version)
	}
	if p.Status != 0x0 {
		t.Fatalf("After unpack, Status in packet is %d, not equal to the expected value: %d\n", p.Status, 0x0)
	}

	authIsmgExpected := []byte{
		0x6c, 0x0b, 0x84, 0x6e, 0x25, 0xba, 0xb6, 0xda,
		0xa4, 0xed, 0x1c, 0x46, 0x6e, 0x0f, 0x4b, 0xd8,
	}

	authIsmg := []byte(p.AuthIsmg)
	for i := 0; i < len(authIsmg); i++ {
		if authIsmg[i] != authIsmgExpected[i] {
			t.Fatalf("After unpack, authIsmg[%d] is %x, not equal to authIsmgExpected[%d]: %x\n", i, authIsmg[i], i, authIsmgExpected[i])
		}
	}
}
