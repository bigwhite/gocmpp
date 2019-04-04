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
	"errors"
	"net"
	"time"
)

var ErrNotCompleted = errors.New("data not being handled completed")
var ErrRespNotMatch = errors.New("the response is not matched with the request")

// Client stands for one client-side instance, just like a session.
// It may connect to the server, send & recv cmpp packets and terminate the connection.
type Client struct {
	conn *Conn
	typ  Type
}

// New establishes a new cmpp client.
func NewClient(typ Type) *Client {
	return &Client{
		typ: typ,
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
	cli.conn = NewConn(conn, cli.typ)
	defer func() {
		if err != nil {
			if cli.conn != nil {
				cli.conn.Close()
			}
		}
	}()
	cli.conn.SetState(CONN_CONNECTED)

	// Login to the server.
	req := &CmppConnReqPkt{
		SrcAddr: user,
		Secret:  password,
		Version: cli.typ,
	}

	_, err = cli.SendReqPkt(req)
	if err != nil {
		return err
	}

	p, err := cli.conn.RecvAndUnpackPkt(timeout)
	if err != nil {
		return err
	}

	var ok bool
	var status uint8
	if cli.typ == V20 || cli.typ == V21 {
		var rsp *Cmpp2ConnRspPkt
		rsp, ok = p.(*Cmpp2ConnRspPkt)
		if !ok {
			err = ErrRespNotMatch
			return err
		}
		status = rsp.Status
	} else {
		var rsp *Cmpp3ConnRspPkt
		rsp, ok = p.(*Cmpp3ConnRspPkt)
		if !ok {
			err = ErrRespNotMatch
			return err
		}
		status = uint8(rsp.Status)
	}

	if status != 0 {
		if status <= ErrnoConnOthers { //ErrnoConnOthers = 5
			err = ConnRspStatusErrMap[status]
		} else {
			err = ConnRspStatusErrMap[ErrnoConnOthers]
		}
		return err
	}

	cli.conn.SetState(CONN_AUTHOK)
	return nil
}

func (cli *Client) Disconnect() {
	if cli.conn != nil {
		cli.conn.Close()
	}
}

// SendReqPkt pack the cmpp request packet structure and send it to the other peer.
func (cli *Client) SendReqPkt(packet Packer) (uint32, error) {
	seq := <-cli.conn.SeqId
	return seq, cli.conn.SendPkt(packet, seq)
}

// SendRspPkt pack the cmpp response packet structure and send it to the other peer.
func (cli *Client) SendRspPkt(packet Packer, seqId uint32) error {
	return cli.conn.SendPkt(packet, seqId)
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (cli *Client) RecvAndUnpackPkt(timeout time.Duration) (interface{}, error) {
	return cli.conn.RecvAndUnpackPkt(timeout)
}
