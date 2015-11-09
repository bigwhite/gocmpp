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

package cmpppacket

import "encoding/binary"

// Packet length const for cmpp active test request and response packets.
const (
	CmppActiveTestReqPktLen uint32 = 12     //12d, 0xc
	CmppActiveTestRspPktLen uint32 = 12 + 1 //13d, 0xd
)

type CmppActiveTestReqPkt struct {
	// session info
	SeqId uint32
}
type CmppActiveTestRspPkt struct {
	Reserved uint8
	// session info
	SeqId uint32
}

// Pack packs the CmppActiveTestReqPkt to bytes stream for client side.
func (p *CmppActiveTestReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmppActiveTestReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_ACTIVE_TEST)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmppActiveTestReqPkt struct.
func (p *CmppActiveTestReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	return r.Error()
}

// Pack packs the CmppActiveTestRspPkt to bytes stream for client side.
func (p *CmppActiveTestRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmppActiveTestRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_ACTIVE_TEST_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	w.WriteByte(p.Reserved)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppActiveTestRspPkt variable.
// After unpack, you will get all value of fields in
// CmppActiveTestRspPkt struct.
func (p *CmppActiveTestRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	p.Reserved = r.ReadByte()
	return r.Error()
}
