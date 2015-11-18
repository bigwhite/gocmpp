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
	"errors"
	"net"
	"time"

	cmppconn "github.com/bigwhite/gocmpp/conn"
	cmpppacket "github.com/bigwhite/gocmpp/packet"
)

var ErrNotCompleted = errors.New("data not being handled completed")
var ErrRespNotMatch = errors.New("the response is not matched with the request")
var ErrConnClosed = errors.New("the conn is closed")

// Client stands for one client-side instance, just like a session.
// It may connect to the server, send & recv cmpp packets and terminate the connection.
type Client struct {
	conn *cmppconn.Conn
	typ  cmpppacket.Type
}

// New establishes a new cmpp client.
func New(typ cmpppacket.Type) *Client {
	return &Client{
		typ: typ,
	}
}

func (cli *Client) Free() {
	if cli != nil {
		if cli.conn != nil {
			cli.conn.Close()
		}
		cli = nil
	}
}

// Connect connect to the cmpp server in block mode.
// It sends login packet, receive and parse connect response packet.
func (cli *Client) Connect(servAddr, user, password string, timeout time.Duration) error {
	var err error
	conn, err := net.DialTimeout("tcp", servAddr, timeout)
	if err != nil {
		return err
	}
	cli.conn = cmppconn.New(conn, cli.typ)
	defer func() {
		if err != nil {
			cli.conn.Close()
		}
	}()

	// Login to the server.
	req := &cmpppacket.CmppConnReqPkt{
		SrcAddr: user,
		Secret:  password,
		Version: cli.typ,
	}

	err = cli.SendReqPkt(req)
	if err != nil {
		return err
	}

	p, err := cli.conn.RecvAndUnpackPkt()
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
		err = ErrRespNotMatch
		return err
	}

	if status != 0 {
		err = cmpppacket.ConnRspStatusErrMap[status]
		return err
	}

	cli.conn.SetState(cmppconn.CONN_AUTHOK)
	return nil
}

// SendReqPkt pack the cmpp request packet structure and send it to the other peer.
func (cli *Client) SendReqPkt(packet cmpppacket.Packer) error {
	if cli.conn == nil {
		return ErrConnClosed
	}
	return cli.conn.SendPkt(packet, <-cli.conn.SeqId)
}

// SendRspPkt pack the cmpp response packet structure and send it to the other peer.
func (cli *Client) SendRspPkt(packet cmpppacket.Packer, seqId uint32) error {
	if cli.conn == nil {
		return ErrConnClosed
	}
	return cli.conn.SendPkt(packet, seqId)
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (cli *Client) RecvAndUnpackPkt() (interface{}, error) {
	if cli.conn == nil {
		return nil, ErrConnClosed
	}
	return cli.conn.RecvAndUnpackPkt()
}
