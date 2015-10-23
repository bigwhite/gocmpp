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
	"strconv"
	"time"
)

const (
	ConnectReqPacketLen = 4 + 4 + 4 + 6 + 16 + 1 + 4
	ConnectRspPacketLen = 4 + 4 + 4 + 4 + 16 + 1
)

func GetCurTimeStamp() (string, uint32) {
	s := time.Now().Format("0102150405")
	i, _ := strconv.Atoi(s)
	return s, uint32(i)
}

type ConnectRequestPacket struct {
	SourceAddr          string
	AuthenticatorSource string
	Version             Type
	Timestamp           uint32
	Secret              string
	SeqId               uint32
}

type ConnectResponsePacket struct {
	Status              uint32
	AuthenticatorIsmg   string
	Version             Type
	Secret              string
	AuthenticatorSource string
	SeqId               uint32
}

func (p *ConnectRequestPacket) Pack(seqId uint32) ([]byte, error) {
	buf := make([]byte, ConnectReqPacketLen)
	packBuf := bytes.NewBuffer(buf)

	// pack header
	err := binary.Write(packBuf, binary.BigEndian, ConnectReqPacketLen)
	if err != nil {
		return nil, err
	}
	err = binary.Write(packBuf, binary.BigEndian, CMPP_CONNECT)
	if err != nil {
		return nil, err
	}
	err = binary.Write(packBuf, binary.BigEndian, seqId)
	if err != nil {
		return nil, err
	}

	var ts string
	ts, p.Timestamp = GetCurTimeStamp()

	// pack body
	packBuf.WriteString(p.SourceAddr)

	md5 := md5.Sum(bytes.Join([][]byte{[]byte(p.SourceAddr),
		make([]byte, 9),
		[]byte(p.Secret),
		[]byte(ts)},
		nil))
	p.AuthenticatorSource = string(md5[:])
	packBuf.WriteString(p.AuthenticatorSource)
	binary.Write(packBuf, binary.BigEndian, p.Version)
	binary.Write(packBuf, binary.BigEndian, p.Timestamp)

	return packBuf.Bytes(), nil
}

func (p *ConnectRequestPacket) Unpack(data []byte) error {
	return nil
}

func (p *ConnectResponsePacket) Pack(seqId uint32) ([]byte, error) {
	buf := make([]byte, ConnectRspPacketLen)
	packBuf := bytes.NewBuffer(buf)

	// pack header
	err := binary.Write(packBuf, binary.BigEndian, ConnectRspPacketLen)
	if err != nil {
		return nil, err
	}
	err = binary.Write(packBuf, binary.BigEndian, CMPP_CONNECT_RESP)
	if err != nil {
		return nil, err
	}
	err = binary.Write(packBuf, binary.BigEndian, seqId)
	if err != nil {
		return nil, err
	}

	// pack body
	binary.Write(packBuf, binary.BigEndian, p.Status)
	var statusBuf bytes.Buffer
	err = binary.Write(&statusBuf, binary.BigEndian, p.Status)
	if err != nil {
		return nil, err
	}
	md5 := md5.Sum(bytes.Join([][]byte{statusBuf.Bytes(),
		[]byte(p.AuthenticatorSource),
		[]byte(p.Secret)},
		nil))
	p.AuthenticatorIsmg = string(md5[:])
	packBuf.WriteString(p.AuthenticatorIsmg)
	binary.Write(packBuf, binary.BigEndian, p.Version)

	return packBuf.Bytes(), nil
}

func (p *ConnectResponsePacket) Unpack(data []byte) error {
	return nil
}
