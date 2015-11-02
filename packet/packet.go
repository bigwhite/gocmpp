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
	"encoding/binary"
	"errors"
	"fmt"
)

type Type int8

const (
	V30 Type = 0x30
	V21 Type = 0x21
	V20 Type = 0x20
)

func (t Type) String() string {
	switch {
	case t == V30:
		return "cmpp30"
	case t == V21:
		return "cmpp21"
	case t == V20:
		return "cmpp20"
	default:
		return "unknown"
	}
}

const (
	CMPP2_PACKET_MAX = 2477
	CMPP2_PACKET_MIN = 12
	CMPP3_PACKET_MAX = 3335
	CMPP3_PACKET_MIN = 12
)

// Protocol errors.
var ErrTotalLengthInvalid = errors.New("total_length in Packet data is invalid")
var ErrCommandIdInvalid = errors.New("command_Id in Packet data is invalid")

type CommandId uint32

const (
	CMPP_REQUEST_MIN, CMPP_RESPONSE_MIN CommandId = iota, 0x80000000 + iota
	CMPP_CONNECT, CMPP_CONNECT_RESP
	CMPP_TERMINATE, CMPP_TERMINATE_RESP
	_, _
	CMPP_SUBMIT, CMPP_SUBMIT_RESP
	CMPP_DELIVER, CMPP_DELIVER_RESP
	CMPP_QUERY, CMPP_QUERY_RESP
	CMPP_CANCEL, CMPP_CANCEL_RESP
	CMPP_ACTIVE_TEST, CMPP_ACTIVE_TEST_RESP
	CMPP_FWD, CMPP_FWD_RESP
	CMPP_MT_ROUTE, CMPP_MT_ROUTE_RESP CommandId = 0x00000010 - 10 + iota, 0x80000010 - 10 + iota
	CMPP_MO_ROUTE, CMPP_MO_ROUTE_RESP
	CMPP_GET_MT_ROUTE, CMPP_GET_MT_ROUTE_RESP
	CMPP_MT_ROUTE_UPDATE, CMPP_MT_ROUTE_UPDATE_RESP
	CMPP_MO_ROUTE_UPDATE, CMPP_MO_ROUTE_UPDATE_RESP
	CMPP_PUSH_MT_ROUTE_UPDATE, CMPP_PUSH_MT_ROUTE_UPDATE_RESP
	CMPP_PUSH_MO_ROUTE_UPDATE, CMPP_PUSH_MO_ROUTE_UPDATE_RESP
	CMPP_GET_MO_ROUTE, CMPP_GET_MO_ROUTE_RESP
	CMPP_REQUEST_MAX, CMPP_RESPONSE_MAX
)

func (id CommandId) String() string {
	if id <= CMPP_FWD && id > CMPP_REQUEST_MIN {
		return []string{
			"CMPP_CONNECT",
			"CMPP_TERMINATE",
			"CMPP_UNKNOWN",
			"CMPP_SUBMIT",
			"CMPP_DELIVER",
			"CMPP_QUERY",
			"CMPP_CANCEL",
			"CMPP_ACTIVE_TEST",
			"CMPP_FWD",
		}[id-1]
	} else if id < CMPP_REQUEST_MAX {
		return []string{
			"CMPP_MT_ROUTE",
			"CMPP_MO_ROUTE",
			"CMPP_GET_MT_ROUTE",
			"CMPP_MT_ROUTE_UPDATE",
			"CMPP_MO_ROUTE_UPDATE",
			"CMPP_PUSH_MT_ROUTE_UPDATE",
			"CMPP_PUSH_MO_ROUTE_UPDATE",
			"CMPP_GET_MO_ROUTE",
		}[id-0x00000010]
	}

	if id < CMPP_FWD_RESP && id > CMPP_RESPONSE_MIN {
		return []string{
			"CMPP_CONNECT_RESP",
			"CMPP_TERMINATE_RESP",
			"CMPP_UNKNOWN",
			"CMPP_SUBMIT_RESP",
			"CMPP_DELIVER_RESP",
			"CMPP_QUERY_RESP",
			"CMPP_CANCEL_RESP",
			"CMPP_ACTIVE_TEST_RESP",
			"CMPP_FWD_RESP",
		}[id-0x80000001]
	} else if id < CMPP_RESPONSE_MAX {
		return []string{
			"CMPP_MT_ROUTE_RESP",
			"CMPP_MO_ROUTE_RESP",
			"CMPP_GET_MT_ROUTE_RESP",
			"CMPP_MT_ROUTE_UPDATE_RESP",
			"CMPP_MO_ROUTE_UPDATE_RESP",
			"CMPP_PUSH_MT_ROUTE_UPDATE_RESP",
			"CMPP_PUSH_MO_ROUTE_UPDATE_RESP",
			"CMPP_GET_MO_ROUTE_RESP",
		}[id-0x80000010]
	}
	return "unknown"
}

type Packer interface {
	Pack(seqId uint32) ([]byte, error)
	Unpack(data []byte) error
}

// OpError is the error type usually returned by functions in the cmpppacket
// package. It describes the operation and the error which the operation caused.
type OpError struct {
	// err is the error that occurred during the operation.
	// it is the origin error.
	err error

	// op is the operation which caused the error, such as
	// some "read" or "write" in packetWriter or packetReader.
	op string
}

func NewOpError(e error, op string) *OpError {
	return &OpError{
		err: e,
		op:  op,
	}
}

func (e *OpError) Error() string {
	if e.err == nil {
		return "<nil>"
	}
	return e.op + " error: " + e.err.Error()
}

func (e *OpError) Cause() error {
	return e.err
}

func (e *OpError) Op() string {
	return e.op
}

type packetWriter struct {
	wb  *bytes.Buffer
	err *OpError
}

func newPacketWriter() *packetWriter {
	return &packetWriter{
		wb: new(bytes.Buffer),
	}
}

func (w *packetWriter) Bytes() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}
	return w.wb.Bytes(), nil
}

func (w *packetWriter) WriteString(s string) {
	if w.err != nil {
		return
	}

	l1 := len(s)
	l2 := l1
	if l2 > 10 {
		l2 = 10
	}

	n, err := w.wb.WriteString(s)
	if err != nil {
		w.err = NewOpError(err,
			fmt.Sprintf("packetWriter.WriteString writes: %s", s[0:l2]))
		return
	}

	if n != l1 {
		w.err = NewOpError(fmt.Errorf("WriteString writes %d bytes, not equal to %d we expected", n, l1),
			fmt.Sprintf("packetWriter.WriteString writes: %s", s[0:l2]))
		return
	}
}

func (w *packetWriter) WriteInt(order binary.ByteOrder, data interface{}) {
	if w.err != nil {
		return
	}

	err := binary.Write(w.wb, order, data)
	if err != nil {
		w.err = NewOpError(err,
			fmt.Sprintf("packetWriter.WriteInt writes: %#v", data))
		return
	}
}

type packetReader struct {
	rb  *bytes.Buffer
	err *OpError
}

func newPacketReader(data []byte) *packetReader {
	return &packetReader{
		rb: bytes.NewBuffer(data),
	}
}

func (r *packetReader) ReadInt(order binary.ByteOrder, value interface{}) {
	if r.err != nil {
		return
	}

	err := binary.Read(r.rb, order, value)
	if err != nil {
		r.err = NewOpError(err,
			"packetReader.ReadInt")
		return
	}
}

func (r *packetReader) ReadBytes(s []byte) {
	if r.err != nil {
		return
	}

	n, err := r.rb.Read(s)
	if err != nil {
		r.err = NewOpError(err,
			"packetReader.ReadBytes")
		return
	}

	if n != len(s) {
		r.err = NewOpError(fmt.Errorf("ReadBytes reads %d bytes, not equal to %d we expected", n, len(s)),
			"packetWriter.ReadBytes")
		return
	}
}

func (r *packetReader) Error() error {
	if r.err != nil {
		return r.err
	}
	return nil
}
