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

package cmpppacket_test

import (
	"testing"

	"github.com/bigwhite/gocmpp/packet"
)

func TestCmppReceiptPktPack(t *testing.T) {
	p := &cmpppacket.CmppReceiptPkt{
		MsgId:          13025908756704198656,
		Stat:           "DELIVRD",
		SubmitTime:     "1511120955",
		DoneTime:       "1511120957",
		DestTerminalId: "13412340000",
		SmscSequence:   0x12345678,
	}

	data, err := p.Pack()
	if err != nil {
		t.Fatal("CmppReceiptPkt pack error:", err)
	}

	dataExpected := []byte{
		0xb4, 0xc5, 0x53, 0x00, 0x00, 0x01, 0x00, 0x00, 0x44, 0x45, 0x4c, 0x49, 0x56, 0x52, 0x44, 0x31,
		0x35, 0x31, 0x31, 0x31, 0x32, 0x30, 0x39, 0x35, 0x35, 0x31, 0x35, 0x31, 0x31, 0x31, 0x32, 0x30,
		0x39, 0x35, 0x37, 0x31, 0x33, 0x34, 0x31, 0x32, 0x33, 0x34, 0x30, 0x30, 0x30, 0x30, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x56, 0x78,
	}

	l1 := len(data)
	l2 := len(dataExpected)
	if l1 != l2 {
		t.Fatalf("After pack, data length is %d, not equal to length expected: %d\n", l1, l2)
	}

	for i := 0; i < l1; i++ {
		if data[i] != dataExpected[i] {
			t.Fatalf("After pack, data[%d] is %x, not equal to dataExpected[%d]: %x\n", i, data[i], i, dataExpected[i])
		}
	}
}

func TestCmppReceiptPktUnpack(t *testing.T) {
	data := []byte{
		0xb4, 0xc5, 0x53, 0x00, 0x00, 0x01, 0x00, 0x00, 0x44, 0x45, 0x4c, 0x49, 0x56, 0x52, 0x44, 0x31,
		0x35, 0x31, 0x31, 0x31, 0x32, 0x30, 0x39, 0x35, 0x35, 0x31, 0x35, 0x31, 0x31, 0x31, 0x32, 0x30,
		0x39, 0x35, 0x37, 0x31, 0x33, 0x34, 0x31, 0x32, 0x33, 0x34, 0x30, 0x30, 0x30, 0x30, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x56, 0x78,
	}

	p := &cmpppacket.CmppReceiptPkt{}
	err := p.Unpack(data)
	if err != nil {
		t.Fatal("CmppReceiptPkt unpack error:", err)
	}

	var resultSet = []struct {
		name          string
		value         interface{}
		expectedValue interface{}
	}{
		{"MsgId", p.MsgId, uint64(13025908756704198656)},
		{"Stat", p.Stat, "DELIVRD"},
		{"SubmitTime", p.SubmitTime, "1511120955"},
		{"DoneTime", p.DoneTime, "1511120957"},
		{"DestTerminalId", p.DestTerminalId, "13412340000"},
		{"SmscSequence", p.SmscSequence, uint32(0x12345678)},
	}

	for _, r := range resultSet {
		if r.value != r.expectedValue {
			t.Fatalf("After unpack, %s in packet is %#v, not equal to the expected value: %#v\n", r.name, r.value, r.expectedValue)
		}
	}
}
