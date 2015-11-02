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
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"strconv"
	"time"

	"github.com/bigwhite/gocmpp/utils"
)

// Packet length const for cmpp connect request and response packets.
const (
	CmppConnReqPktLen  uint32 = 4 + 4 + 4 + 6 + 16 + 1 + 4 //39d, 0x27
	Cmpp2ConnRspPktLen uint32 = 4 + 4 + 4 + 1 + 16 + 1     //30d, 0x1e
	Cmpp3ConnRspPktLen uint32 = 4 + 4 + 4 + 4 + 16 + 1     //33d, 0x21
)

// Errors for connect resp status.
var ErrConnInvalidStruct = errors.New("connect response status: invalid protocol structure")
var ErrConnInvalidSrcAddr = errors.New("connect response status: invalid source address")
var ErrConnAuthFailed = errors.New("connect response status: Auth failed")
var ErrConnVerTooHigh = errors.New("connect response status: protocol version is too high")
var ErrConnOthers = errors.New("connect response status: other errors")

var ConnRspStatusErrMap = map[uint8]error{
	1: ErrConnInvalidStruct,
	2: ErrConnInvalidSrcAddr,
	3: ErrConnAuthFailed,
	4: ErrConnVerTooHigh,
	5: ErrConnOthers,
}

func now() (string, uint32) {
	s := time.Now().Format("0102150405")
	i, _ := strconv.Atoi(s)
	return s, uint32(i)
}

// CmppConnReqPkt represents a Cmpp2 or Cmpp3 connect request packet.
//
// when used in client side(pack), you should initialize it with
// correct SourceAddr(SrcAddr), Secret and Version.
//
// when used in server side(unpack), nothing needed to be initialized.
// unpack will fill the SourceAddr(SrcAddr), AuthSrc, Version, Timestamp
// and SeqId
//
type CmppConnReqPkt struct {
	SrcAddr   string
	AuthSrc   string
	Version   Type
	Timestamp uint32
	Secret    string
	SeqId     uint32
}

// Cmpp2ConnRspPkt represents a Cmpp2 connect response packet.
//
// when used in server side(pack), you should initialize it with
// correct Status, AuthSrc, Secret and Version.
//
// when used in client side(unpack), nothing needed to be initialized.
// unpack will fill the Status, AuthImsg, Version and SeqId
//
type Cmpp2ConnRspPkt struct {
	Status   uint8
	AuthIsmg string
	Version  Type
	Secret   string
	AuthSrc  string
	SeqId    uint32
}

// Cmpp3ConnRspPkt represents a Cmpp3 connect response packet.
//
// when used in server side(pack), you should initialize it with
// correct Status, AuthSrc, Secret and Version.
//
// when used in client side(unpack), nothing needed to be initialized.
// unpack will fill the Status, AuthImsg, Version and SeqId
//
type Cmpp3ConnRspPkt struct {
	Status   uint32
	AuthIsmg string
	Version  Type
	Secret   string
	AuthSrc  string
	SeqId    uint32
}

// Pack packs the CmppConnReqPkt to bytes stream for client side.
// Before calling Pack, you should initialize a CmppConnReqPkt variable
// with correct SourceAddr(SrcAddr), Secret and Version.
func (p *CmppConnReqPkt) Pack(seqId uint32) ([]byte, error) {
	var w = newPacketWriter()

	// Pack header
	w.WriteInt(binary.BigEndian, CmppConnReqPktLen)
	w.WriteInt(binary.BigEndian, CMPP_CONNECT)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	var ts string
	if p.Timestamp == 0 {
		ts, p.Timestamp = now() //default: current time.
	} else {
		ts = cmpputils.TimeStamp2Str(p.Timestamp)
	}

	// Pack body
	w.WriteString(p.SrcAddr)

	md5 := md5.Sum(bytes.Join([][]byte{[]byte(p.SrcAddr),
		make([]byte, 9),
		[]byte(p.Secret),
		[]byte(ts)},
		nil))
	p.AuthSrc = string(md5[:])

	w.WriteString(p.AuthSrc)
	w.WriteInt(binary.BigEndian, p.Version)
	w.WriteInt(binary.BigEndian, p.Timestamp)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppConnReqPkt variable.
// Usually it is used in server side. After unpack, you will get SeqId, SourceAddr,
// AuthenticatorSource, Version and Timestamp.
func (p *CmppConnReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	// Body: Source_Addr
	var sa = make([]byte, 6)
	r.ReadBytes(sa)
	p.SrcAddr = string(sa)

	// Body: AuthSrc
	var as = make([]byte, 16)
	r.ReadBytes(as)
	p.AuthSrc = string(as)

	// Body: Version
	r.ReadInt(binary.BigEndian, &p.Version)
	// Body: timestamp
	r.ReadInt(binary.BigEndian, &p.Timestamp)

	return r.Error()
}

// Pack packs the Cmpp2ConnRspPkt to bytes stream for server side.
// Before calling Pack, you should initialize a Cmpp2ConnRspPkt variable
// with correct Status,AuthenticatorSource, Secret and Version.
func (p *Cmpp2ConnRspPkt) Pack(seqId uint32) ([]byte, error) {
	var w = newPacketWriter()

	// pack header
	w.WriteInt(binary.BigEndian, Cmpp2ConnRspPktLen)
	w.WriteInt(binary.BigEndian, CMPP_CONNECT_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// pack body
	w.WriteInt(binary.BigEndian, p.Status)

	md5 := md5.Sum(bytes.Join([][]byte{[]byte{p.Status},
		[]byte(p.AuthSrc),
		[]byte(p.Secret)},
		nil))
	p.AuthIsmg = string(md5[:])
	w.WriteString(p.AuthIsmg)

	w.WriteInt(binary.BigEndian, p.Version)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp2ConnRspPkt variable.
// Usually it is used in client side. After unpack, you will get SeqId, Status,
// AuthenticatorIsmg, and Version.
// Parameter data contains seqId in header and the whole packet body.
func (p *Cmpp2ConnRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	// Body: Status
	r.ReadInt(binary.BigEndian, &p.Status)

	// Body: AuthenticatorISMG
	var s = make([]byte, 16)
	r.ReadBytes(s)
	p.AuthIsmg = string(s)

	// Body: Version
	r.ReadInt(binary.BigEndian, &p.Version)
	return r.Error()
}

// Pack packs the Cmpp3ConnRspPkt to bytes stream for server side.
// Before calling Pack, you should initialize a Cmpp3ConnRspPkt variable
// with correct Status,AuthenticatorSource, Secret and Version.
func (p *Cmpp3ConnRspPkt) Pack(seqId uint32) ([]byte, error) {
	var w = newPacketWriter()

	// pack header
	w.WriteInt(binary.BigEndian, Cmpp3ConnRspPktLen)
	w.WriteInt(binary.BigEndian, CMPP_CONNECT_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	// pack body
	w.WriteInt(binary.BigEndian, p.Status)

	var statusBuf = new(bytes.Buffer)
	err := binary.Write(statusBuf, binary.BigEndian, p.Status)
	if err != nil {
		return nil, err
	}

	md5 := md5.Sum(bytes.Join([][]byte{statusBuf.Bytes(),
		[]byte(p.AuthSrc),
		[]byte(p.Secret)},
		nil))
	p.AuthIsmg = string(md5[:])
	w.WriteString(p.AuthIsmg)

	w.WriteInt(binary.BigEndian, p.Version)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3ConnRspPkt variable.
// Usually it is used in client side. After unpack, you will get SeqId, Status,
// AuthenticatorIsmg, and Version.
// Parameter data contains seqId in header and the whole packet body.
func (p *Cmpp3ConnRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)

	// Body: Status
	r.ReadInt(binary.BigEndian, &p.Status)

	// Body: AuthenticatorISMG
	var s = make([]byte, 16)
	r.ReadBytes(s)
	p.AuthIsmg = string(s)

	// Body: Version
	r.ReadInt(binary.BigEndian, &p.Version)
	return r.Error()
}
