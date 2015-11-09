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

type Cmpp2DeliverReqPkt struct {
}
type Cmpp2DeliverRspPkt struct {
}
type Cmpp3DeliverReqPkt struct {
}
type Cmpp3DeliverRspPkt struct {
}

// Pack packs the Cmpp2DeliverReqPkt to bytes stream for client side.
func (p *Cmpp2DeliverReqPkt) Pack(seqId uint32) ([]byte, error) {
	return nil, nil
}

// Unpack unpack the binary byte stream to a Cmpp2DeliverReqPkt variable.
// After unpack, you will get all value of fields in
// Cmpp2DeliverReqPkt struct.
func (p *Cmpp2DeliverReqPkt) Unpack(data []byte) error {
	return nil
}

// Pack packs the Cmpp2DeliverRspPkt to bytes stream for client side.
func (p *Cmpp2DeliverRspPkt) Pack(seqId uint32) ([]byte, error) {
	return nil, nil
}

// Unpack unpack the binary byte stream to a Cmpp2DeliverRspPkt variable.
// After unpack, you will get all value of fields in
// Cmpp2DeliverRspPkt struct.
func (p *Cmpp2DeliverRspPkt) Unpack(data []byte) error {
	return nil
}

// Pack packs the Cmpp3DeliverReqPkt to bytes stream for client side.
func (p *Cmpp3DeliverReqPkt) Pack(seqId uint32) ([]byte, error) {
	return nil, nil
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverReqPkt variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverReqPkt struct.
func (p *Cmpp3DeliverReqPkt) Unpack(data []byte) error {
	return nil
}

// Pack packs the Cmpp3DeliverRspPkt to bytes stream for client side.
func (p *Cmpp3DeliverRspPkt) Pack(seqId uint32) ([]byte, error) {
	return nil, nil
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverRspPkt variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverRspPkt struct.
func (p *Cmpp3DeliverRspPkt) Unpack(data []byte) error {
	return nil
}
