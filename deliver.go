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

import (
	"encoding/binary"
	"errors"
)

// Packet length const for cmpp deliver request and response packets.
const (
	Cmpp2DeliverReqPktMaxLen uint32 = 12 + 233   //245d, 0xf5
	Cmpp2DeliverRspPktLen    uint32 = 12 + 8 + 1 //21d, 0x15

	Cmpp3DeliverReqPktMaxLen uint32 = 12 + 257   //269d, 0x10d
	Cmpp3DeliverRspPktLen    uint32 = 12 + 8 + 4 //24d, 0x18
)

// Errors for result in deliver resp.

var (
	ErrnoDeliverInvalidStruct      uint8 = 1
	ErrnoDeliverInvalidCommandId   uint8 = 2
	ErrnoDeliverInvalidSequence    uint8 = 3
	ErrnoDeliverInvalidMsgLength   uint8 = 4
	ErrnoDeliverInvalidFeeCode     uint8 = 5
	ErrnoDeliverExceedMaxMsgLength uint8 = 6
	ErrnoDeliverInvalidServiceId   uint8 = 7
	ErrnoDeliverNotPassFlowControl uint8 = 8
	ErrnoDeliverOtherError         uint8 = 9

	DeliverRspResultErrMap = map[uint8]error{
		ErrnoDeliverInvalidStruct:      errDeliverInvalidStruct,
		ErrnoDeliverInvalidCommandId:   errDeliverInvalidCommandId,
		ErrnoDeliverInvalidSequence:    errDeliverInvalidSequence,
		ErrnoDeliverInvalidMsgLength:   errDeliverInvalidMsgLength,
		ErrnoDeliverInvalidFeeCode:     errDeliverInvalidFeeCode,
		ErrnoDeliverExceedMaxMsgLength: errDeliverExceedMaxMsgLength,
		ErrnoDeliverInvalidServiceId:   errDeliverInvalidServiceId,
		ErrnoDeliverNotPassFlowControl: errDeliverNotPassFlowControl,
		ErrnoDeliverOtherError:         errDeliverOtherError,
	}

	errDeliverInvalidStruct      = errors.New("deliver response status: invalid protocol structure")
	errDeliverInvalidCommandId   = errors.New("deliver response status: invalid command id")
	errDeliverInvalidSequence    = errors.New("deliver response status: invalid message sequence")
	errDeliverInvalidMsgLength   = errors.New("deliver response status: invalid message length")
	errDeliverInvalidFeeCode     = errors.New("deliver response status: invalid fee code")
	errDeliverExceedMaxMsgLength = errors.New("deliver response status: exceed max message length")
	errDeliverInvalidServiceId   = errors.New("deliver response status: invalid service id")
	errDeliverNotPassFlowControl = errors.New("deliver response status: not pass the flow control")
	errDeliverOtherError         = errors.New("deliver response status: other error")
)

type Cmpp2DeliverReqPkt struct {
	MsgId            uint64
	DestId           string
	ServiceId        string
	TpPid            uint8
	TpUdhi           uint8
	MsgFmt           uint8
	SrcTerminalId    string
	RegisterDelivery uint8
	MsgLength        uint8
	MsgContent       string
	Reserve          string

	//session info
	SeqId uint32
}

type Cmpp2DeliverRspPkt struct {
	MsgId  uint64
	Result uint8

	//session info
	SeqId uint32
}
type Cmpp3DeliverReqPkt struct {
	MsgId            uint64
	DestId           string
	ServiceId        string
	TpPid            uint8
	TpUdhi           uint8
	MsgFmt           uint8
	SrcTerminalId    string
	SrcTerminalType  uint8
	RegisterDelivery uint8
	MsgLength        uint8
	MsgContent       string
	LinkId           string

	//session info
	SeqId uint32
}
type Cmpp3DeliverRspPkt struct {
	MsgId  uint64
	Result uint32

	//session info
	SeqId uint32
}

// Pack packs the Cmpp2DeliverReqPkt to bytes stream for client side.
func (p *Cmpp2DeliverReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CMPP_HEADER_LEN + 65 + uint32(p.MsgLength) + 8

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_DELIVER)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteFixedSizeString(p.DestId, 21)
	w.WriteFixedSizeString(p.ServiceId, 10)
	w.WriteByte(p.TpPid)
	w.WriteByte(p.TpUdhi)
	w.WriteByte(p.MsgFmt)
	w.WriteFixedSizeString(p.SrcTerminalId, 21)
	w.WriteByte(p.RegisterDelivery)
	w.WriteByte(p.MsgLength)
	w.WriteString(p.MsgContent)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp2DeliverReqPkt variable.
// After unpack, you will get all value of fields in
// Cmpp2DeliverReqPkt struct.
func (p *Cmpp2DeliverReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	// Body
	r.ReadInt(binary.BigEndian, &p.MsgId)

	destId := r.ReadCString(21)
	p.DestId = string(destId)

	serviceId := r.ReadCString(10)
	p.ServiceId = string(serviceId)

	p.TpPid = r.ReadByte()
	p.TpUdhi = r.ReadByte()
	p.MsgFmt = r.ReadByte()

	srcTerminalId := r.ReadCString(21)
	p.SrcTerminalId = string(srcTerminalId)

	p.RegisterDelivery = r.ReadByte()
	p.MsgLength = r.ReadByte()

	msgContent := make([]byte, p.MsgLength)
	r.ReadBytes(msgContent)
	p.MsgContent = string(msgContent)

	reserve := r.ReadCString(8)
	p.Reserve = string(reserve)

	return r.Error()
}

// Pack packs the Cmpp2DeliverRspPkt to bytes stream for client side.
func (p *Cmpp2DeliverRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = Cmpp2DeliverRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_DELIVER_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteByte(p.Result)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp2DeliverRspPkt variable.
// After unpack, you will get all value of fields in
// Cmpp2DeliverRspPkt struct.
func (p *Cmpp2DeliverRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	r.ReadInt(binary.BigEndian, &p.MsgId)
	p.Result = r.ReadByte()

	return r.Error()
}

// Pack packs the Cmpp3DeliverReqPkt to bytes stream for client side.
func (p *Cmpp3DeliverReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CMPP_HEADER_LEN + 77 + uint32(p.MsgLength) + 20

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_DELIVER)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteFixedSizeString(p.DestId, 21)
	w.WriteFixedSizeString(p.ServiceId, 10)
	w.WriteByte(p.TpPid)
	w.WriteByte(p.TpUdhi)
	w.WriteByte(p.MsgFmt)
	w.WriteFixedSizeString(p.SrcTerminalId, 32)
	w.WriteByte(p.SrcTerminalType)
	w.WriteByte(p.RegisterDelivery)
	w.WriteByte(p.MsgLength)
	w.WriteString(p.MsgContent)
	w.WriteFixedSizeString(p.LinkId, 20)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverReqPkt variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverReqPkt struct.
func (p *Cmpp3DeliverReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	// Body
	r.ReadInt(binary.BigEndian, &p.MsgId)

	destId := r.ReadCString(21)
	p.DestId = string(destId)

	serviceId := r.ReadCString(10)
	p.ServiceId = string(serviceId)

	p.TpPid = r.ReadByte()
	p.TpUdhi = r.ReadByte()
	p.MsgFmt = r.ReadByte()

	srcTerminalId := r.ReadCString(32)
	p.SrcTerminalId = string(srcTerminalId)
	p.SrcTerminalType = r.ReadByte()

	p.RegisterDelivery = r.ReadByte()
	p.MsgLength = r.ReadByte()

	msgContent := make([]byte, p.MsgLength)
	r.ReadBytes(msgContent)
	p.MsgContent = string(msgContent)

	linkId := r.ReadCString(20)
	p.LinkId = string(linkId)

	return r.Error()
}

// Pack packs the Cmpp3DeliverRspPkt to bytes stream for client side.
func (p *Cmpp3DeliverRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = Cmpp3DeliverRspPktLen
	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_DELIVER_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteInt(binary.BigEndian, p.Result)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverRspPkt variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverRspPkt struct.
func (p *Cmpp3DeliverRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	r.ReadInt(binary.BigEndian, &p.MsgId)
	r.ReadInt(binary.BigEndian, &p.Result)

	return r.Error()
}
