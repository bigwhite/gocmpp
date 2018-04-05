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
	"io"
	"net"
	"sync"
	"time"
)

type State uint8

// Errors for conn operations
var (
	ErrConnIsClosed       = errors.New("connection is closed")
	ErrReadCmdIDTimeout   = errors.New("read commandId timeout")
	ErrReadPktBodyTimeout = errors.New("read packet body timeout")
)

var noDeadline = time.Time{}

// Conn States
const (
	CONN_CLOSED State = iota
	CONN_CONNECTED
	CONN_AUTHOK
)

type Conn struct {
	net.Conn
	State State
	Typ   Type

	// for SeqId generator goroutine
	SeqId <-chan uint32
	done  chan<- struct{}
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

// New returns an abstract structure for successfully
// established underlying net.Conn.
func NewConn(conn net.Conn, typ Type) *Conn {
	seqId, done := newSeqIdGenerator()
	c := &Conn{
		Conn:  conn,
		Typ:   typ,
		SeqId: seqId,
		done:  done,
	}
	tc := c.Conn.(*net.TCPConn) // Always tcpconn
	tc.SetKeepAlive(true)       //Keepalive as default
	return c
}

func (c *Conn) Close() {
	if c != nil {
		if c.State == CONN_CLOSED {
			return
		}
		close(c.done)  // let the SeqId goroutine exit.
		c.Conn.Close() // close the underlying net.Conn
		c.State = CONN_CLOSED
	}
}

func (c *Conn) SetState(state State) {
	c.State = state
}

// SendPkt pack the cmpp packet structure and send it to the other peer.
func (c *Conn) SendPkt(packet Packer, seqId uint32) error {
	if c.State == CONN_CLOSED {
		return ErrConnIsClosed
	}

	data, err := packet.Pack(seqId)
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(data) //block write
	if err != nil {
		return err
	}

	return nil
}

const (
	defaultReadBufferSize = 4096
)

// readBuffer is used to optimize the performance of
// RecvAndUnpackPkt.
type readBuffer struct {
	totalLen  uint32
	commandId CommandId
	leftData  [defaultReadBufferSize]byte
}

var readBufferPool = sync.Pool{
	New: func() interface{} {
		return &readBuffer{}
	},
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (c *Conn) RecvAndUnpackPkt(timeout time.Duration) (interface{}, error) {
	if c.State == CONN_CLOSED {
		return nil, ErrConnIsClosed
	}
	defer c.SetReadDeadline(noDeadline)

	rb := readBufferPool.Get().(*readBuffer)
	defer readBufferPool.Put(rb)

	// Total_Length in packet
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
	err := binary.Read(c.Conn, binary.BigEndian, &rb.totalLen)
	if err != nil {
		return nil, err
	}

	if c.Typ == V30 {
		if rb.totalLen < CMPP3_PACKET_MIN || rb.totalLen > CMPP3_PACKET_MAX {
			return nil, ErrTotalLengthInvalid
		}
	}

	if c.Typ == V21 || c.Typ == V20 {
		if rb.totalLen < CMPP2_PACKET_MIN || rb.totalLen > CMPP2_PACKET_MAX {
			return nil, ErrTotalLengthInvalid
		}
	}

	// Command_Id
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
	err = binary.Read(c.Conn, binary.BigEndian, &rb.commandId)
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok {
			if netErr.Timeout() {
				return nil, ErrReadCmdIDTimeout
			}
		}
		return nil, err
	}

	if !((rb.commandId > CMPP_REQUEST_MIN && rb.commandId < CMPP_REQUEST_MAX) ||
		(rb.commandId > CMPP_RESPONSE_MIN && rb.commandId < CMPP_RESPONSE_MAX)) {
		return nil, ErrCommandIdInvalid
	}

	// The left packet data (start from seqId in header).
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
	var leftData = rb.leftData[0:(rb.totalLen - 8)]
	_, err = io.ReadFull(c.Conn, leftData)
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok {
			if netErr.Timeout() {
				return nil, ErrReadPktBodyTimeout
			}
		}
		return nil, err
	}

	var p Packer
	switch rb.commandId {
	case CMPP_CONNECT:
		p = &CmppConnReqPkt{}
	case CMPP_CONNECT_RESP:
		if c.Typ == V30 {
			p = &Cmpp3ConnRspPkt{}
		} else {
			p = &Cmpp2ConnRspPkt{}
		}
	case CMPP_TERMINATE:
		p = &CmppTerminateReqPkt{}
	case CMPP_TERMINATE_RESP:
		p = &CmppTerminateRspPkt{}
	case CMPP_SUBMIT:
		if c.Typ == V30 {
			p = &Cmpp3SubmitReqPkt{}
		} else {
			p = &Cmpp2SubmitReqPkt{}
		}
	case CMPP_SUBMIT_RESP:
		if c.Typ == V30 {
			p = &Cmpp3SubmitRspPkt{}
		} else {
			p = &Cmpp2SubmitRspPkt{}
		}
	case CMPP_DELIVER:
		if c.Typ == V30 {
			p = &Cmpp3DeliverReqPkt{}
		} else {
			p = &Cmpp2DeliverReqPkt{}
		}
	case CMPP_DELIVER_RESP:
		if c.Typ == V30 {
			p = &Cmpp3DeliverRspPkt{}
		} else {
			p = &Cmpp2DeliverRspPkt{}
		}
	case CMPP_FWD:
		if c.Typ == V30 {
			p = &Cmpp3FwdReqPkt{}
		} else {
			p = &Cmpp2FwdReqPkt{}
		}
	case CMPP_FWD_RESP:
		if c.Typ == V30 {
			p = &Cmpp3FwdRspPkt{}
		} else {
			p = &Cmpp2FwdRspPkt{}
		}
	case CMPP_ACTIVE_TEST:
		p = &CmppActiveTestReqPkt{}
	case CMPP_ACTIVE_TEST_RESP:
		p = &CmppActiveTestRspPkt{}

	default:
		p = nil
		return nil, ErrCommandIdNotSupported
	}

	err = p.Unpack(leftData)
	if err != nil {
		return nil, err
	}
	return p, nil
}
