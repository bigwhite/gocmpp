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

// Packet length const for cmpp receipt packet.
const (
	CmppReceiptPktLen uint32 = 60 //60d, 0x3c
)

type CmppReceiptPkt struct {
	MsgId          uint64
	Stat           string
	SubmitTime     string // YYMMDDHHMM
	DoneTime       string // YYMMDDHHMM
	DestTerminalId string
	SmscSequence   uint32
}

// Pack packs the CmppReceiptPkt to bytes stream for client side.
func (p *CmppReceiptPkt) Pack() ([]byte, error) {
	var pktLen uint32 = CmppReceiptPktLen

	var w = newPacketWriter(pktLen)

	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteFixedSizeString(p.Stat, 7)
	w.WriteFixedSizeString(p.SubmitTime, 10)
	w.WriteFixedSizeString(p.DoneTime, 10)
	w.WriteFixedSizeString(p.DestTerminalId, 21)
	w.WriteInt(binary.BigEndian, p.SmscSequence)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppReceiptPkt variable.
// After unpack, you will get all value of fields in
// CmppReceiptPkt struct.
func (p *CmppReceiptPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	r.ReadInt(binary.BigEndian, &p.MsgId)

	stat := r.ReadCString(7)
	p.Stat = string(stat)

	submitTime := r.ReadCString(10)
	p.SubmitTime = string(submitTime)

	doneTime := r.ReadCString(10)
	p.DoneTime = string(doneTime)

	destTerminalId := r.ReadCString(21)
	p.DestTerminalId = string(destTerminalId)

	r.ReadInt(binary.BigEndian, &p.SmscSequence)
	return r.Error()
}
