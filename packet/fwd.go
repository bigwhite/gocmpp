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

type Cmpp2FwdReqPkt struct {
}
type Cmpp2FwdRspPkt struct {
}
type Cmpp3FwdReqPkt struct {
}
type Cmpp3FwdRspPkt struct {
}

func (p *Cmpp3FwdReqPkt) Pack(seqId uint32) ([]byte, error) {
	return nil, nil
}
func (p *Cmpp3FwdReqPkt) Unpack(data []byte) error {
	return nil
}

func (p *Cmpp3FwdRspPkt) Pack(seqId uint32) ([]byte, error) {
	return nil, nil

}
func (p *Cmpp3FwdRspPkt) Unpack(data []byte) error {
	return nil
}

func (p *Cmpp2FwdReqPkt) Pack(seqId uint32) ([]byte, error) {
	return nil, nil
}
func (p *Cmpp2FwdReqPkt) Unpack(data []byte) error {
	return nil
}

func (p *Cmpp2FwdRspPkt) Pack(seqId uint32) ([]byte, error) {
	return nil, nil
}
func (p *Cmpp2FwdRspPkt) Unpack(data []byte) error {
	return nil
}
