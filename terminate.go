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

package cmpp

import "encoding/binary"

// Packet length const for cmpp terminate request and response packets.
const (
	CmppTerminateReqPktLen uint32 = 12 //12d, 0xc
	CmppTerminateRspPktLen uint32 = 12 //12d, 0xc
)

type CmppTerminateReqPkt struct {
	// session info
	SeqId uint32
}
type CmppTerminateRspPkt struct {
	// session info
	SeqId uint32
}

// Pack packs the CmppTerminateReqPkt to bytes stream for client side.
func (p *CmppTerminateReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmppTerminateReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_TERMINATE)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateReqPkt variable.
// After unpack, you will get all value of fields in
// CmppTerminateReqPkt struct.
func (p *CmppTerminateReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	return r.Error()
}

// Pack packs the CmppTerminateRspPkt to bytes stream for client side.
func (p *CmppTerminateRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmppTerminateRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_TERMINATE_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateRspPkt variable.
// After unpack, you will get all value of fields in
// CmppTerminateRspPkt struct.
func (p *CmppTerminateRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	return r.Error()
}
