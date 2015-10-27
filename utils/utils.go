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

package cmpputils

import (
	"fmt"
	"unsafe"
)

func IsBigEndian() bool {
	var i uint16 = 0x1234
	var p *[2]byte = (*[2]byte)(unsafe.Pointer(&i))
	if (*p)[0] == 0x12 {
		return true
	}
	return false
}

// TimeStamp2Str converts a timestamp(MMDDHHMMSS) int to a string(10 bytes).
func TimeStamp2Str(t uint32) string {
	return fmt.Sprintf("%010d", t)
}
