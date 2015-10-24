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

import "errors"

type Type int8

const (
	Ver30 Type = 0x30
	Ver21 Type = 0x21
	Ver20 Type = 0x20
)

func (t Type) String() string {
	switch {
	case t == Ver30:
		return "cmpp30"
	case t == Ver21:
		return "cmpp21"
	case t == Ver20:
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

// Protocol errors
var ErrTotalLengthInvalid = errors.New("Total_length in Packet data is invalid")
var ErrCommandIdInvalid = errors.New("Command_Id in Packet data is invalid")

type CommandId uint32

const (
	CMPP_REQUEST_MIN, CMPP_RESPONSE_MIN CommandId = iota, 0x80000000 + iota
	CMPP_CONNECT, CMPP_CONNECT_RESP
	CMPP_TERMINATE, CMPP_TERMINATE_RESP
	_, _
	CMPP_SUBMIT, CMPP_SUBMIT_RESP
	CMPP_DELIVER, CMPP_DELIVER_RESP
	CMPP_REQUEST_MAX, CMPP_RESPONSE_MAX
)

func (id CommandId) String() string {
	if id < CMPP_REQUEST_MAX && id > CMPP_REQUEST_MIN {
		return []string{
			"CMPP_CONNECT",
			"CMPP_TERMINATE",
		}[id-1]
	}
	if id < CMPP_RESPONSE_MAX && id > CMPP_RESPONSE_MIN {
		return []string{
			"CMPP_CONNECT_RESP",
			"CMPP_TERMINATE_RESP",
		}[id-0x80000001]
	}
	return "unknown"
}

type Packer interface {
	Pack(seqId uint32) ([]byte, error)
	Unpack(data []byte) error
}
