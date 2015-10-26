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
	"strings"
	"time"
)

const (
	ConnectReqPacketLen uint32 = 4 + 4 + 4 + 6 + 16 + 1 + 4
	ConnectRspPacketLen uint32 = 4 + 4 + 4 + 4 + 16 + 1
)

// Errors for connect resp status
var ErrConnInvalidStruct = errors.New("Connect response status: invalid protocol structure")
var ErrConnInvalidSourceAddr = errors.New("Connect response status: invalid source address")
var ErrConnAuthFailed = errors.New("Connect response status: Auth failed")
var ErrConnVerTooHigh = errors.New("Connect response status: protocol version is too high")
var ErrConnOthers = errors.New("Connect response status: other errors")

var ConnRespStatusErrMap = map[uint8]error{
	1: ErrConnInvalidStruct,
	2: ErrConnInvalidSourceAddr,
	3: ErrConnAuthFailed,
	4: ErrConnVerTooHigh,
	5: ErrConnOthers,
}

func now() (string, uint32) {
	s := time.Now().Format("0102150405")
	i, _ := strconv.Atoi(s)
	return s, uint32(i)
}

// timeStamp2Str converts a timestamp(MMDDHHMMSS) int to a string(10 bytes).
func timeStamp2Str(t uint32) string {
	s := strconv.Itoa(int(t))
	n := 10 - len(s)

	if n == 0 {
		return s
	} else if n > 0 {
		var buf = make([]byte, n)
		for i := 0; i < n; i++ {
			buf[i] = '0'
		}
		return strings.Join([]string{string(buf), s}, "")
	}
	return "" //should never reach here.
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
	Status              uint8
	AuthenticatorIsmg   string
	Version             Type
	Secret              string
	AuthenticatorSource string
	SeqId               uint32
}

func (p *ConnectRequestPacket) Pack(seqId uint32) ([]byte, error) {
	var packBuf = new(bytes.Buffer)

	// pack header
	err := binary.Write(packBuf, binary.BigEndian, ConnectReqPacketLen)
	if err != nil {
		return nil, err
	}
	err = binary.Write(packBuf, binary.BigEndian, uint32(CMPP_CONNECT))
	if err != nil {
		return nil, err
	}
	err = binary.Write(packBuf, binary.BigEndian, seqId)
	if err != nil {
		return nil, err
	}
	p.SeqId = seqId

	var ts string
	if p.Timestamp == 0 {
		ts, p.Timestamp = now()
	} else {
		ts = timeStamp2Str(p.Timestamp)
	}

	// pack body
	_, err = packBuf.WriteString(p.SourceAddr)
	if err != nil {
		return nil, err
	}

	md5 := md5.Sum(bytes.Join([][]byte{[]byte(p.SourceAddr),
		make([]byte, 9),
		[]byte(p.Secret),
		[]byte(ts)},
		nil))
	p.AuthenticatorSource = string(md5[:])

	_, err = packBuf.WriteString(p.AuthenticatorSource)
	if err != nil {
		return nil, err
	}

	err = binary.Write(packBuf, binary.BigEndian, p.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Write(packBuf, binary.BigEndian, p.Timestamp)
	if err != nil {
		return nil, err
	}

	return packBuf.Bytes(), nil
}

func (p *ConnectRequestPacket) Unpack(data []byte) error {
	var buf = bytes.NewBuffer(data)

	// Sequence Id
	err := binary.Read(buf, binary.BigEndian, &p.SeqId)
	if err != nil {
		return err
	}

	// Body: Source_Addr
	var sa = make([]byte, 6)
	_, err = buf.Read(sa)
	if err != nil {
		return err
	}
	p.SourceAddr = string(sa)

	// Body: AuthenticatorSource
	var as = make([]byte, 16)
	_, err = buf.Read(as)
	if err != nil {
		return err
	}
	p.AuthenticatorSource = string(as)

	// Body: Version
	err = binary.Read(buf, binary.BigEndian, &p.Version)
	if err != nil {
		return err
	}

	// Body: timestamp
	err = binary.Read(buf, binary.BigEndian, &p.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (p *ConnectResponsePacket) Pack(seqId uint32) ([]byte, error) {
	var packBuf = new(bytes.Buffer)

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
	err = binary.Write(packBuf, binary.BigEndian, p.Status)
	if err != nil {
		return nil, err
	}

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

	_, err = packBuf.WriteString(p.AuthenticatorIsmg)
	if err != nil {
		return nil, err
	}

	err = binary.Write(packBuf, binary.BigEndian, p.Version)
	if err != nil {
		return nil, err
	}

	return packBuf.Bytes(), nil
}

// data include seqId in header + the whole body
func (p *ConnectResponsePacket) Unpack(data []byte) error {
	var buf = bytes.NewBuffer(data)

	// Sequence Id
	err := binary.Read(buf, binary.BigEndian, &p.SeqId)
	if err != nil {
		return err
	}

	// Body: Status
	err = binary.Read(buf, binary.BigEndian, &p.Status)
	if err != nil {
		return err
	}

	// Body: AuthenticatorISMG
	var s = make([]byte, 16)
	_, err = buf.Read(s)
	if err != nil {
		return err
	}
	p.AuthenticatorIsmg = string(s)

	// Body: Version
	err = binary.Read(buf, binary.BigEndian, &p.Version)
	if err != nil {
		return err
	}
	return nil
}
