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

import (
	"encoding/binary"
	"errors"
)

// Packet length const for cmpp fwd request and response packets.
const (
	Cmpp2FwdReqPktMaxLen uint32 = 12 + 2379          //2277d, 0x957
	Cmpp2FwdRspPktLen    uint32 = 12 + 8 + 1 + 1 + 1 //23d, 0x17

	Cmpp3FwdReqPktMaxLen uint32 = 12 + 2491          //2503d, 0x9c7
	Cmpp3FwdRspPktLen    uint32 = 12 + 8 + 1 + 1 + 4 //26d, 0x1a
)

// Errors for result in fwd resp.
var ErrFwdInvalidStruct = errors.New("fwd response status: invalid protocol structure")
var ErrFwdInvalidCommandId = errors.New("fwd response status: invalid command id")
var ErrFwdInvalidSequence = errors.New("fwd response status: invalid message sequence")
var ErrFwdInvalidMsgLength = errors.New("fwd response status: invalid message length")
var ErrFwdInvalidFeeCode = errors.New("fwd response status: invalid fee code")
var ErrFwdExceedMaxMsgLength = errors.New("fwd response status: exceed max message length")
var ErrFwdInvalidServiceId = errors.New("fwd response status: invalid service id")
var ErrFwdNotPassFlowControl = errors.New("fwd response status: not pass the flow control")
var ErrFwdNoPrivilege = errors.New("fwd response status: msg has no fwd privilege")

var FwdRspResultErrMap = map[uint8]error{
	1: ErrFwdInvalidStruct,
	2: ErrFwdInvalidCommandId,
	3: ErrFwdInvalidSequence,
	4: ErrFwdInvalidMsgLength,
	5: ErrFwdInvalidFeeCode,
	6: ErrFwdExceedMaxMsgLength,
	7: ErrFwdInvalidServiceId,
	8: ErrFwdNotPassFlowControl,
	9: ErrFwdNoPrivilege,
}

type Cmpp2FwdReqPkt struct {
	SourceId           string
	DestinationId      string
	NodesCount         uint8
	MsgFwdType         uint8
	MsgId              uint64
	PkTotal            uint8
	PkNumber           uint8
	RegisteredDelivery uint8
	MsgLevel           uint8
	ServiceId          string
	FeeUserType        uint8
	FeeTerminalId      string
	TpPid              uint8
	TpUdhi             uint8
	MsgFmt             uint8
	MsgSrc             string
	FeeType            string
	FeeCode            string
	ValidTime          string
	AtTime             string
	SrcId              string
	DestUsrTl          uint8
	DestId             []string
	MsgLength          uint8
	MsgContent         string

	// session info
	SeqId uint32
}

type Cmpp2FwdRspPkt struct {
	MsgId    uint64
	PkTotal  uint8
	PkNumber uint8
	Result   uint8

	// session info
	SeqId uint32
}
type Cmpp3FwdReqPkt struct {
	SourceId            string
	DestinationId       string
	NodesCount          uint8
	MsgFwdType          uint8
	MsgId               uint64
	PkTotal             uint8
	PkNumber            uint8
	RegisteredDelivery  uint8
	MsgLevel            uint8
	ServiceId           string
	FeeUserType         uint8
	FeeTerminalId       string
	FeeTerminalPseudo   string
	FeeTerminalUserType uint8
	TpPid               uint8
	TpUdhi              uint8
	MsgFmt              uint8
	MsgSrc              string
	FeeType             string
	FeeCode             string
	ValidTime           string
	AtTime              string
	SrcId               string
	SrcPseudo           string
	SrcUserType         uint8
	SrcType             uint8
	DestUsrTl           uint8
	DestId              []string
	DestPseudo          string
	DestUserType        uint8
	MsgLength           uint8
	MsgContent          string
	LinkId              string

	// session info
	SeqId uint32
}

type Cmpp3FwdRspPkt struct {
	MsgId    uint64
	PkTotal  uint8
	PkNumber uint8
	Result   uint32

	// session info
	SeqId uint32
}

// Pack packs the Cmpp2FwdReqPkt to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp2FwdReqPkt variable
// with correct field value.
func (p *Cmpp2FwdReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CMPP_HEADER_LEN + 131 + uint32(p.DestUsrTl)*21 + 1 + uint32(p.MsgLength) + 8
	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_FWD)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteFixedSizeString(p.SourceId, 6)
	w.WriteFixedSizeString(p.DestinationId, 6)
	w.WriteByte(p.NodesCount)
	w.WriteByte(p.MsgFwdType)
	w.WriteInt(binary.BigEndian, p.MsgId)

	if p.PkTotal == 0 && p.PkNumber == 0 {
		p.PkTotal, p.PkNumber = 1, 1
	}
	w.WriteByte(p.PkTotal)
	w.WriteByte(p.PkNumber)

	w.WriteByte(p.RegisteredDelivery)
	w.WriteByte(p.MsgLevel)
	w.WriteFixedSizeString(p.ServiceId, 10)
	w.WriteByte(p.FeeUserType)
	w.WriteFixedSizeString(p.FeeTerminalId, 21)
	w.WriteByte(p.TpPid)
	w.WriteByte(p.TpUdhi)
	w.WriteByte(p.MsgFmt)
	w.WriteFixedSizeString(p.MsgSrc, 6)
	w.WriteString(p.FeeType)
	w.WriteFixedSizeString(p.FeeCode, 6)
	w.WriteFixedSizeString(p.ValidTime, 17)
	w.WriteFixedSizeString(p.AtTime, 17)
	w.WriteFixedSizeString(p.SrcId, 21)
	w.WriteByte(p.DestUsrTl)
	for _, d := range p.DestId {
		w.WriteFixedSizeString(d, 21)
	}
	w.WriteByte(p.MsgLength)
	w.WriteString(p.MsgContent)
	w.WriteFixedSizeString("", 8) //Reserved

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp2FwdReqPkt variable.
// After unpack, you will get all value of fields in Cmpp2FwdReqPkt struct.
func (p *Cmpp2FwdReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	sourceId := r.ReadCString(6)
	p.SourceId = string(sourceId)
	destinationId := r.ReadCString(6)
	p.DestinationId = string(destinationId)
	p.NodesCount = r.ReadByte()
	p.MsgFwdType = r.ReadByte()

	r.ReadInt(binary.BigEndian, &p.MsgId)

	p.PkTotal = r.ReadByte()
	p.PkNumber = r.ReadByte()
	p.RegisteredDelivery = r.ReadByte()
	p.MsgLevel = r.ReadByte()
	serviceId := r.ReadCString(10)
	p.ServiceId = string(serviceId)
	p.FeeUserType = r.ReadByte()
	feeTerminalId := r.ReadCString(21)
	p.FeeTerminalId = string(feeTerminalId)
	p.TpPid = r.ReadByte()
	p.TpUdhi = r.ReadByte()
	p.MsgFmt = r.ReadByte()

	msgSrc := r.ReadCString(6)
	p.MsgSrc = string(msgSrc)

	feeType := make([]byte, 2)
	r.ReadBytes(feeType)
	p.FeeType = string(feeType)

	feeCode := r.ReadCString(6)
	p.FeeCode = string(feeCode)

	validTime := r.ReadCString(17)
	p.ValidTime = string(validTime)

	atTime := r.ReadCString(17)
	p.AtTime = string(atTime)

	srcId := r.ReadCString(21)
	p.SrcId = string(srcId)

	p.DestUsrTl = r.ReadByte()
	for i := 0; i < int(p.DestUsrTl); i++ {
		destId := r.ReadCString(21)
		p.DestId = append(p.DestId, string(destId))
	}

	p.MsgLength = r.ReadByte()

	msgContent := make([]byte, p.MsgLength)
	r.ReadBytes(msgContent)
	p.MsgContent = string(msgContent)

	return r.Error()
}

// Pack packs the Cmpp2FwdRspPkt to bytes stream for server side.
// Before calling Pack, you should initialize a Cmpp2FwdRspPkt variable
// with correct field value.
func (p *Cmpp2FwdRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = Cmpp2FwdRspPktLen
	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_FWD_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteByte(p.PkTotal)
	w.WriteByte(p.PkNumber)
	w.WriteByte(p.Result)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp2FwdRspPkt variable.
// After unpack, you will get all value of fields in Cmpp2FwdRspPkt struct.
func (p *Cmpp2FwdRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	r.ReadInt(binary.BigEndian, &p.MsgId)
	p.PkTotal = r.ReadByte()
	p.PkNumber = r.ReadByte()
	p.Result = r.ReadByte()

	return r.Error()
}

// Pack packs the Cmpp3FwdReqPkt to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp3FwdReqPkt variable
// with correct field value.
func (p *Cmpp3FwdReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CMPP_HEADER_LEN + 198 + uint32(p.DestUsrTl)*21 + 32 + 1 + 1 + uint32(p.MsgLength) + 20

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_FWD)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteFixedSizeString(p.SourceId, 6)
	w.WriteFixedSizeString(p.DestinationId, 6)
	w.WriteByte(p.NodesCount)
	w.WriteByte(p.MsgFwdType)
	w.WriteInt(binary.BigEndian, p.MsgId)

	if p.PkTotal == 0 && p.PkNumber == 0 {
		p.PkTotal, p.PkNumber = 1, 1
	}
	w.WriteByte(p.PkTotal)
	w.WriteByte(p.PkNumber)
	w.WriteByte(p.RegisteredDelivery)
	w.WriteByte(p.MsgLevel)
	w.WriteFixedSizeString(p.ServiceId, 10)
	w.WriteByte(p.FeeUserType)
	w.WriteFixedSizeString(p.FeeTerminalId, 21)
	w.WriteFixedSizeString(p.FeeTerminalPseudo, 32)
	w.WriteByte(p.FeeTerminalUserType)
	w.WriteByte(p.TpPid)
	w.WriteByte(p.TpUdhi)
	w.WriteByte(p.MsgFmt)
	w.WriteFixedSizeString(p.MsgSrc, 6)
	w.WriteString(p.FeeType)
	w.WriteFixedSizeString(p.FeeCode, 6)
	w.WriteFixedSizeString(p.ValidTime, 17)
	w.WriteFixedSizeString(p.AtTime, 17)
	w.WriteFixedSizeString(p.SrcId, 21)
	w.WriteFixedSizeString(p.SrcPseudo, 32)
	w.WriteByte(p.SrcUserType)
	w.WriteByte(p.SrcType)
	w.WriteByte(p.DestUsrTl)

	for _, d := range p.DestId {
		w.WriteFixedSizeString(d, 21)
	}
	w.WriteFixedSizeString(p.DestPseudo, 32)
	w.WriteByte(p.DestUserType)
	w.WriteByte(p.MsgLength)
	w.WriteString(p.MsgContent)
	w.WriteFixedSizeString(p.LinkId, 20)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3FwdReqPkt variable.
// After unpack, you will get all value of fields in Cmpp3FwdReqPkt struct.
func (p *Cmpp3FwdReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	// Body
	sourceId := r.ReadCString(6)
	p.SourceId = string(sourceId)
	destinationId := r.ReadCString(6)
	p.DestinationId = string(destinationId)
	p.NodesCount = r.ReadByte()
	p.MsgFwdType = r.ReadByte()

	r.ReadInt(binary.BigEndian, &p.MsgId)

	p.PkTotal = r.ReadByte()
	p.PkNumber = r.ReadByte()
	p.RegisteredDelivery = r.ReadByte()
	p.MsgLevel = r.ReadByte()

	serviceId := r.ReadCString(10)
	p.ServiceId = string(serviceId)

	p.FeeUserType = r.ReadByte()

	feeTerminalId := r.ReadCString(21)
	p.FeeTerminalId = string(feeTerminalId)
	feeTerminalPseudo := r.ReadCString(32)
	p.FeeTerminalPseudo = string(feeTerminalPseudo)
	p.FeeTerminalUserType = r.ReadByte()

	p.TpPid = r.ReadByte()
	p.TpUdhi = r.ReadByte()
	p.MsgFmt = r.ReadByte()

	msgSrc := r.ReadCString(6)
	p.MsgSrc = string(msgSrc)

	feeType := make([]byte, 2)
	r.ReadBytes(feeType)
	p.FeeType = string(feeType)

	feeCode := r.ReadCString(6)
	p.FeeCode = string(feeCode)

	validTime := r.ReadCString(17)
	p.ValidTime = string(validTime)

	atTime := r.ReadCString(17)
	p.AtTime = string(atTime)

	srcId := r.ReadCString(21)
	p.SrcId = string(srcId)

	srcPseudo := r.ReadCString(32)
	p.SrcPseudo = string(srcPseudo)
	p.SrcUserType = r.ReadByte()
	p.SrcType = r.ReadByte()

	p.DestUsrTl = r.ReadByte()
	for i := 0; i < int(p.DestUsrTl); i++ {
		destId := r.ReadCString(21)
		p.DestId = append(p.DestId, string(destId))
	}
	destPseudo := r.ReadCString(32)
	p.DestPseudo = string(destPseudo)
	p.DestUserType = r.ReadByte()

	p.MsgLength = r.ReadByte()
	msgContent := make([]byte, p.MsgLength)
	r.ReadBytes(msgContent)
	p.MsgContent = string(msgContent)

	linkId := r.ReadCString(20)
	p.LinkId = string(linkId)

	return r.Error()
}

// Pack packs the Cmpp3FwdRspPkt to bytes stream for server side.
// Before calling Pack, you should initialize a Cmpp3FwdRspPkt variable
// with correct field value.
func (p *Cmpp3FwdRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = Cmpp3FwdRspPktLen
	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_FWD_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteByte(p.PkTotal)
	w.WriteByte(p.PkNumber)
	w.WriteInt(binary.BigEndian, p.Result)

	return w.Bytes()

}

// Unpack unpack the binary byte stream to a Cmpp3FwdRspPkt variable.
// After unpack, you will get all value of fields in Cmpp3FwdRspPkt struct.
func (p *Cmpp3FwdRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	r.ReadInt(binary.BigEndian, &p.MsgId)
	p.PkTotal = r.ReadByte()
	p.PkNumber = r.ReadByte()
	r.ReadInt(binary.BigEndian, &p.Result)

	return r.Error()
}
