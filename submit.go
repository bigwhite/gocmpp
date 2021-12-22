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

// Packet length const for cmpp submit request and response packets.
const (
	Cmpp2SubmitReqPktMaxLen uint32 = 12 + 2265  //2277d, 0x8e5
	Cmpp2SubmitRspPktLen    uint32 = 12 + 8 + 1 //21d, 0x15

	Cmpp3SubmitReqPktMaxLen uint32 = 12 + 3479  //3491d, 0xda3
	Cmpp3SubmitRspPktLen    uint32 = 12 + 8 + 4 //24d, 0x18
)

// Errors for result in submit resp.
var (
	ErrnoSubmitInvalidStruct         uint8 = 1
	ErrnoSubmitInvalidCommandId      uint8 = 2
	ErrnoSubmitInvalidSequence       uint8 = 3
	ErrnoSubmitInvalidMsgLength      uint8 = 4
	ErrnoSubmitInvalidFeeCode        uint8 = 5
	ErrnoSubmitExceedMaxMsgLength    uint8 = 6
	ErrnoSubmitInvalidServiceId      uint8 = 7
	ErrnoSubmitNotPassFlowControl    uint8 = 8
	ErrnoSubmitNotServeFeeTerminalId uint8 = 9
	ErrnoSubmitInvalidSrcId          uint8 = 10
	ErrnoSubmitInvalidMsgSrc         uint8 = 11
	ErrnoSubmitInvalidFeeTerminalId  uint8 = 12
	ErrnoSubmitInvalidDestTerminalId uint8 = 13

	SubmitRspResultErrMap = map[uint8]error{
		ErrnoSubmitInvalidStruct:         errSubmitInvalidStruct,
		ErrnoSubmitInvalidCommandId:      errSubmitInvalidCommandId,
		ErrnoSubmitInvalidSequence:       errSubmitInvalidSequence,
		ErrnoSubmitInvalidMsgLength:      errSubmitInvalidMsgLength,
		ErrnoSubmitInvalidFeeCode:        errSubmitInvalidFeeCode,
		ErrnoSubmitExceedMaxMsgLength:    errSubmitExceedMaxMsgLength,
		ErrnoSubmitInvalidServiceId:      errSubmitInvalidServiceId,
		ErrnoSubmitNotPassFlowControl:    errSubmitNotPassFlowControl,
		ErrnoSubmitNotServeFeeTerminalId: errSubmitNotServeFeeTerminalId,
		ErrnoSubmitInvalidSrcId:          errSubmitInvalidSrcId,
		ErrnoSubmitInvalidMsgSrc:         errSubmitInvalidMsgSrc,
		ErrnoSubmitInvalidFeeTerminalId:  errSubmitInvalidFeeTerminalId,
		ErrnoSubmitInvalidDestTerminalId: errSubmitInvalidDestTerminalId,
	}

	errSubmitInvalidStruct         = errors.New("submit response status: invalid protocol structure")
	errSubmitInvalidCommandId      = errors.New("submit response status: invalid command id")
	errSubmitInvalidSequence       = errors.New("submit response status: invalid message sequence")
	errSubmitInvalidMsgLength      = errors.New("submit response status: invalid message length")
	errSubmitInvalidFeeCode        = errors.New("submit response status: invalid fee code")
	errSubmitExceedMaxMsgLength    = errors.New("submit response status: exceed max message length")
	errSubmitInvalidServiceId      = errors.New("submit response status: invalid service id")
	errSubmitNotPassFlowControl    = errors.New("submit response status: not pass the flow control")
	errSubmitNotServeFeeTerminalId = errors.New("submit response status: feeTerminalId is not served")
	errSubmitInvalidSrcId          = errors.New("submit response status: invalid srcId")
	errSubmitInvalidMsgSrc         = errors.New("submit response status: invalid msgSrc")
	errSubmitInvalidFeeTerminalId  = errors.New("submit response status: invalid feeTerminalId")
	errSubmitInvalidDestTerminalId = errors.New("submit response status: invalid destTerminalId")
)

type Cmpp2SubmitReqPkt struct {
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
	DestTerminalId     []string
	MsgLength          uint8
	MsgContent         string
	Reserve            string

	// session info
	SeqId uint32
}

type Cmpp2SubmitRspPkt struct {
	MsgId  uint64
	Result uint8

	// session info
	SeqId uint32
}

type Cmpp3SubmitReqPkt struct {
	MsgId              uint64
	PkTotal            uint8
	PkNumber           uint8
	RegisteredDelivery uint8
	MsgLevel           uint8
	ServiceId          string
	FeeUserType        uint8
	FeeTerminalId      string
	FeeTerminalType    uint8
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
	DestTerminalId     []string
	DestTerminalType   uint8
	MsgLength          uint8
	MsgContent         string
	LinkId             string

	// session info
	SeqId uint32
}

type Cmpp3SubmitRspPkt struct {
	MsgId  uint64
	Result uint32

	// session info
	SeqId uint32
}

// Pack packs the Cmpp2SubmitReqPkt to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp2SubmitReqPkt variable
// with correct field value.
func (p *Cmpp2SubmitReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CMPP_HEADER_LEN + 117 + uint32(p.DestUsrTl)*21 + 1 + uint32(p.MsgLength) + 8

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_SUBMIT)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
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
	w.WriteFixedSizeString(p.FeeType, 2)
	w.WriteFixedSizeString(p.FeeCode, 6)
	w.WriteFixedSizeString(p.ValidTime, 17)
	w.WriteFixedSizeString(p.AtTime, 17)
	w.WriteFixedSizeString(p.SrcId, 21)
	w.WriteByte(p.DestUsrTl)

	for _, d := range p.DestTerminalId {
		w.WriteFixedSizeString(d, 21)
	}
	w.WriteByte(p.MsgLength)
	w.WriteString(p.MsgContent)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp2SubmitReqPkt variable.
// Usually it is used in server side. After unpack, you will get all value of fields in
// Cmpp2SubmitReqPkt struct.
func (p *Cmpp2SubmitReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

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
		destTerminalId := r.ReadCString(21)
		p.DestTerminalId = append(p.DestTerminalId, string(destTerminalId))
	}

	p.MsgLength = r.ReadByte()

	msgContent := make([]byte, p.MsgLength)
	r.ReadBytes(msgContent)
	p.MsgContent = string(msgContent)

	reserve := r.ReadCString(8)
	p.Reserve = string(reserve)

	return r.Error()
}

// Pack packs the Cmpp2SubmitRspPkt to bytes stream for Server side.
// Before calling Pack, you should initialize a Cmpp2SubmitRspPkt variable
// with correct field value.
func (p *Cmpp2SubmitRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CMPP_HEADER_LEN + 8 + 1

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_SUBMIT_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteByte(p.Result)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp2SubmitRspPkt variable.
// Usually it is used in client side. After unpack, you will get all value of fields in
// Cmpp2SubmitRspPkt struct.
func (p *Cmpp2SubmitRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	r.ReadInt(binary.BigEndian, &p.MsgId)
	p.Result = r.ReadByte()

	return r.Error()
}

// Pack packs the Cmpp3SubmitReqPkt to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp3SubmitReqPkt variable
// with correct field value.
func (p *Cmpp3SubmitReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CMPP_HEADER_LEN + 129 + uint32(p.DestUsrTl)*32 + 1 + 1 + uint32(p.MsgLength) + 20

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_SUBMIT)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
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
	w.WriteFixedSizeString(p.FeeTerminalId, 32)
	w.WriteByte(p.FeeTerminalType)
	w.WriteByte(p.TpPid)
	w.WriteByte(p.TpUdhi)
	w.WriteByte(p.MsgFmt)
	w.WriteFixedSizeString(p.MsgSrc, 6)
	w.WriteFixedSizeString(p.FeeType, 2)
	w.WriteFixedSizeString(p.FeeCode, 6)
	w.WriteFixedSizeString(p.ValidTime, 17)
	w.WriteFixedSizeString(p.AtTime, 17)
	w.WriteFixedSizeString(p.SrcId, 21)
	w.WriteByte(p.DestUsrTl)

	for _, d := range p.DestTerminalId {
		w.WriteFixedSizeString(d, 32)
	}
	w.WriteByte(p.DestTerminalType)
	w.WriteByte(p.MsgLength)
	w.WriteString(p.MsgContent)
	w.WriteFixedSizeString(p.LinkId, 20)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3SubmitReqPkt variable.
// Usually it is used in server side. After unpack, you will get all value of fields in
// Cmpp3SubmitReqPkt struct.
func (p *Cmpp3SubmitReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	r.ReadInt(binary.BigEndian, &p.MsgId)

	p.PkTotal = r.ReadByte()
	p.PkNumber = r.ReadByte()
	p.RegisteredDelivery = r.ReadByte()
	p.MsgLevel = r.ReadByte()

	serviceId := r.ReadCString(10)
	p.ServiceId = string(serviceId)

	p.FeeUserType = r.ReadByte()

	feeTerminalId := r.ReadCString(32)
	p.FeeTerminalId = string(feeTerminalId)

	p.FeeTerminalType = r.ReadByte()
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
		destTerminalId := r.ReadCString(32)
		p.DestTerminalId = append(p.DestTerminalId, string(destTerminalId))
	}

	p.DestTerminalType = r.ReadByte()
	p.MsgLength = r.ReadByte()

	msgContent := make([]byte, p.MsgLength)
	r.ReadBytes(msgContent)
	p.MsgContent = string(msgContent)

	linkId := r.ReadCString(20)
	p.LinkId = string(linkId)

	return r.Error()
}

// Pack packs the Cmpp3SubmitRspPkt to bytes stream for Server side.
// Before calling Pack, you should initialize a Cmpp3SubmitRspPkt variable
// with correct field value.
func (p *Cmpp3SubmitRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CMPP_HEADER_LEN + 8 + 4

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMPP_SUBMIT_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// Pack Body
	w.WriteInt(binary.BigEndian, p.MsgId)
	w.WriteInt(binary.BigEndian, p.Result)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3SubmitRspPkt variable.
// Usually it is used in client side. After unpack, you will get all value of fields in
// Cmpp3SubmitRspPkt struct.
func (p *Cmpp3SubmitRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	r.ReadInt(binary.BigEndian, &p.MsgId)
	r.ReadInt(binary.BigEndian, &p.Result)

	return r.Error()
}
