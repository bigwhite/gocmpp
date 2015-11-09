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

package cmppclient

import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	cmpppacket "github.com/bigwhite/gocmpp/packet"
)

var ErrNotCompleted = errors.New("data not being handled completed")
var ErrRespNotMatch = errors.New("the response is not matched with the request")

// Client stands for one client-side instance, just like a session.
// It may connect to the server, send & recv cmpp packets and terminate the connection.
type Client struct {
	t         uint8 // packet response timeout, default: 60s
	keepAlive bool  // indicates whether current session is a keepalive one, default: true
	conn      net.Conn
	typ       cmpppacket.Type
	seqId     <-chan uint32
	done      chan<- struct{}
}

// New establishes a new cmpp client.
func New(typ cmpppacket.Type) *Client {
	seqId, done := newSeqIdGenerator()
	return &Client{
		t:         60,
		keepAlive: true,
		typ:       typ,
		seqId:     seqId,
		done:      done,
	}
}

func newSeqIdGenerator() (<-chan uint32, chan<- struct{}) {
	out := make(chan uint32)
	done := make(chan struct{})

	go func() {
		var i uint32
		for {
			select {
			case out <- i:
				i++
			case <-done:
				close(out)
				return
			}
		}
	}()
	return out, done
}

func (cli *Client) Free() {
	if cli != nil {
		if cli.conn != nil {
			cli.conn.Close()
		}
		close(cli.done)
		cli = nil
	}
}

// SetT sets the heartbeat response timeout for the client.
// You should call this method before session established.
func (cli *Client) SetT(t uint8) {
	cli.t = t
}

// SetKeepAlive sets the connection attribute for the client.
// You should call this method before session established.
func (cli *Client) SetKeepAlive(keepAlive bool) {
	cli.keepAlive = keepAlive
}

// Connect connect to the cmpp server in block mode.
// It sends login packet, receive and parse connect response packet.
func (cli *Client) Connect(servAddr, user, password string) error {
	var err error
	cli.conn, err = net.Dial("tcp", servAddr)
	if err != nil {
		return err
	}

	// Login to the server.
	req := &cmpppacket.CmppConnReqPkt{
		SrcAddr: user,
		Secret:  password,
		Version: cli.typ,
	}

	err = cli.SendPacket(req)
	if err != nil {
		return err
	}

	p, err := cli.RecvAndUnpackPkt()
	if err != nil {
		return err
	}

	var ok bool
	var status uint8
	if cli.typ == cmpppacket.V20 || cli.typ == cmpppacket.V21 {
		var rsp *cmpppacket.Cmpp2ConnRspPkt
		rsp, ok = p.(*cmpppacket.Cmpp2ConnRspPkt)
		status = rsp.Status
	} else {
		var rsp *cmpppacket.Cmpp3ConnRspPkt
		rsp, ok = p.(*cmpppacket.Cmpp3ConnRspPkt)
		status = uint8(rsp.Status)
	}

	if !ok {
		return ErrRespNotMatch
	}

	if status != 0 {
		return cmpppacket.ConnRspStatusErrMap[status]
	}

	return nil
}

// SendPacket pack the cmpp packet structure and send it to the other peer.
func (cli *Client) SendPacket(packet cmpppacket.Packer) error {
	data, err := packet.Pack(<-cli.seqId)
	if err != nil {
		return err
	}

	n, err := cli.conn.Write(data)
	if err != nil {
		return nil
	}

	if n != len(data) {
		return ErrNotCompleted
	}
	return nil
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (cli *Client) RecvAndUnpackPkt() (interface{}, error) {
	// Total_Length in packet
	var totalLen uint32
	err := binary.Read(cli.conn, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}

	if cli.typ == cmpppacket.V30 {
		if totalLen < cmpppacket.CMPP3_PACKET_MIN || totalLen > cmpppacket.CMPP3_PACKET_MAX {
			return nil, cmpppacket.ErrTotalLengthInvalid
		}
	}

	if cli.typ == cmpppacket.V21 || cli.typ == cmpppacket.V20 {
		if totalLen < cmpppacket.CMPP2_PACKET_MIN || totalLen > cmpppacket.CMPP2_PACKET_MAX {
			return nil, cmpppacket.ErrTotalLengthInvalid
		}
	}

	// Command_Id
	var commandId cmpppacket.CommandId
	err = binary.Read(cli.conn, binary.BigEndian, &commandId)
	if err != nil {
		return nil, err
	}

	if !((commandId > cmpppacket.CMPP_REQUEST_MIN && commandId < cmpppacket.CMPP_REQUEST_MAX) ||
		(commandId > cmpppacket.CMPP_RESPONSE_MIN && commandId < cmpppacket.CMPP_RESPONSE_MAX)) {
		return nil, cmpppacket.ErrCommandIdInvalid
	}

	// The left packet data (start from seqId in header).
	// todo: may use cli.conn.SetReadDeadline to avoid longtime block
	var leftData = make([]byte, totalLen-8)
	_, err = io.ReadFull(cli.conn, leftData)
	if err != nil {
		return nil, err
	}

	var p cmpppacket.Packer
	switch commandId {
	case cmpppacket.CMPP_CONNECT:
		p = &cmpppacket.CmppConnReqPkt{}
	case cmpppacket.CMPP_CONNECT_RESP:
		if cli.typ == cmpppacket.V30 {
			p = &cmpppacket.Cmpp3ConnRspPkt{}
		} else {
			p = &cmpppacket.Cmpp2ConnRspPkt{}
		}
	case cmpppacket.CMPP_TERMINATE:
		p = &cmpppacket.CmppTerminateReqPkt{}
	case cmpppacket.CMPP_TERMINATE_RESP:
		p = &cmpppacket.CmppTerminateRspPkt{}
	case cmpppacket.CMPP_SUBMIT:
		if cli.typ == cmpppacket.V30 {
			p = &cmpppacket.Cmpp3SubmitReqPkt{}
		} else {
			p = &cmpppacket.Cmpp2SubmitReqPkt{}
		}
	case cmpppacket.CMPP_SUBMIT_RESP:
		if cli.typ == cmpppacket.V30 {
			p = &cmpppacket.Cmpp3SubmitRspPkt{}
		} else {
			p = &cmpppacket.Cmpp2SubmitRspPkt{}
		}
	case cmpppacket.CMPP_DELIVER:
		if cli.typ == cmpppacket.V30 {
			p = &cmpppacket.Cmpp3DeliverReqPkt{}
		} else {
			p = &cmpppacket.Cmpp2DeliverReqPkt{}
		}
	case cmpppacket.CMPP_DELIVER_RESP:
		if cli.typ == cmpppacket.V30 {
			p = &cmpppacket.Cmpp3DeliverRspPkt{}
		} else {
			p = &cmpppacket.Cmpp2DeliverRspPkt{}
		}
	case cmpppacket.CMPP_FWD:
		if cli.typ == cmpppacket.V30 {
			p = &cmpppacket.Cmpp3FwdReqPkt{}
		} else {
			p = &cmpppacket.Cmpp2FwdReqPkt{}
		}
	case cmpppacket.CMPP_FWD_RESP:
		if cli.typ == cmpppacket.V30 {
			p = &cmpppacket.Cmpp3FwdRspPkt{}
		} else {
			p = &cmpppacket.Cmpp2FwdRspPkt{}
		}
	case cmpppacket.CMPP_ACTIVE_TEST:
		p = &cmpppacket.CmppActiveTestReqPkt{}
	case cmpppacket.CMPP_ACTIVE_TEST_RESP:
		p = &cmpppacket.CmppActiveTestRspPkt{}

	default:
		p = nil
		return nil, cmpppacket.ErrCommandIdNotSupported
	}

	err = p.Unpack(leftData)
	if err != nil {
		return nil, err
	}
	return p, nil
}
