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

package cmpputils_test

import (
	"testing"

	cmpputils "github.com/bigwhite/gocmpp/utils"
)

func TestTimeStamp2Str(t *testing.T) {
	var t1 uint32 = 1021080510
	s1 := cmpputils.TimeStamp2Str(t1)
	if s1 != "1021080510" {
		t.Errorf("The result of TimeStamp2Str is %s, not equal to expected: %s\n", s1, "1021080510")
	}

	var t2 uint32 = 121080510
	s2 := cmpputils.TimeStamp2Str(t2)
	if s2 != "0121080510" {
		t.Errorf("The result of TimeStamp2Str is %s, not equal to expected: %s\n", s2, "0121080510")
	}
}

func TestUtf8ToUcs2(t *testing.T) {
	// invalid utf8 bytes sequences
	b1 := []byte{0xe6, 0xb1, 0x89, 0xe6}
	_, err := cmpputils.Utf8ToUcs2(string(b1))
	if err == nil {
		t.Fatal("The in string passed to Utf8ToUcs2 is valid, not equal to our expected")
	}

	if err != cmpputils.ErrInvalidUtf8Rune {
		t.Fatalf("The result is %#v, not equal to our expected %#v", err, cmpputils.ErrInvalidUtf8Rune)
	}

	// valid utf8 bytes sequences
	b2 := []byte{0xe6, 0xb1, 0x89}
	s2, err := cmpputils.Utf8ToUcs2(string(b2))
	if err != nil {
		t.Fatalf("The error is %#v, not to the result expected: nil", err)
	}

	sl2 := []byte(s2)

	if sl2[0] != 0x6c {
		t.Fatalf("The first char is %x, not equal to expected %x\n", sl2[0], 0x6c)
	}

	if sl2[1] != 0x49 {
		t.Fatalf("The second char is %x, not equal to expected %x\n", sl2[1], 0x49)
	}
}

func TestUcs2ToUtf8(t *testing.T) {
	u1 := []byte{0x6c, 0x49}

	s1, err := cmpputils.Ucs2ToUtf8(string(u1))
	if err != nil {
		t.Fatalf("The error is %#v, not to the result expected: nil", err)
	}

	if s1 != "汉" {
		t.Fatalf("The result is %s, not equal to our expected %s", s1, "汉")
	}

}

func TestUtf8ToGB18030(t *testing.T) {
	// invalid utf8 bytes sequences
	b1 := []byte{0xe6, 0xb1, 0x89, 0xe6}
	_, err := cmpputils.Utf8ToGB18030(string(b1))
	if err == nil {
		t.Fatal("The in string passed to Utf8ToGB18030is valid, not equal to our expected")
	}

	if err != cmpputils.ErrInvalidUtf8Rune {
		t.Fatalf("The result is %#v, not equal to our expected %#v", err, cmpputils.ErrInvalidUtf8Rune)
	}

	// valid utf8 bytes sequences
	b2 := []byte{0xe6, 0xb1, 0x89} //"汉"
	s2, err := cmpputils.Utf8ToGB18030(string(b2))
	if err != nil {
		t.Fatalf("The error is %#v, not to the result expected: nil", err)
	}

	sl2 := []byte(s2)

	if sl2[0] != 0xba {
		t.Fatalf("The first char is %x, not equal to expected %x\n", sl2[0], 0xba)
	}

	if sl2[1] != 0xba {
		t.Fatalf("The second char is %x, not equal to expected %x\n", sl2[1], 0xba)
	}

}

func TestGB18030ToUtf8(t *testing.T) {
	u1 := []byte{0xd6, 0xd0}

	s1, err := cmpputils.GB18030ToUtf8(string(u1))
	if err != nil {
		t.Fatalf("The error is %#v, not to the result expected: nil", err)
	}

	if s1 != "中" {
		t.Fatalf("The result is %s, not equal to our expected %s", s1, "中")
	}
}

func TestOctetString(t *testing.T) {
	s1 := "666666"
	s2 := "88888888"
	s3 := "55555"

	s := cmpputils.OctetString(s1, 6)
	if s != "666666" {
		t.Fatalf("The result is %s, not equal to our expected %s", s, "666666")
	}

	s = cmpputils.OctetString(s2, 6)
	if s != "888888" {
		t.Fatalf("The result is %s, not equal to our expected %s", s, "888888")
	}

	s = cmpputils.OctetString(s3, 6)
	expected := s3 + string(make([]byte, 1))
	if s != expected {
		t.Fatalf("The result is %s, not equal to our expected %s", s, expected)
	}
}
