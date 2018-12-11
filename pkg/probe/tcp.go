/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package probe

import (
	"encoding/binary"
)

type tcpFlags struct {
	ns, cwr, ece, urg, ack, psh, rst, syn, fin bool
}

type tcpHeader struct {
	srcPort    uint16 // 0
	destPort   uint16 // 2
	seq        uint32 // 4
	ack        uint32 // 8
	dataOffset uint16
	flags      tcpFlags // 12
	windowSize uint16   // 14
	checkSum   uint16   // 16
	urgentPtr  uint16   // 18
}

func (t *tcpHeader) Encode() []byte {
	pkt := make([]byte, 20)

	encoder := binary.BigEndian
	encoder.PutUint16(pkt, t.srcPort)
	encoder.PutUint16(pkt[2:], t.destPort)
	encoder.PutUint32(pkt[4:], t.seq)
	encoder.PutUint32(pkt[8:], t.ack)

	flagAt := func(b bool, offset uint) uint16 {
		if !b {
			return 0
		}
		return 1 << offset
	}

	var flags uint16
	flags |= t.dataOffset & 0x3
	// Reserved (3 bits).
	flags |= flagAt(t.flags.ns, 7)
	flags |= flagAt(t.flags.cwr, 8)
	flags |= flagAt(t.flags.ece, 9)
	flags |= flagAt(t.flags.urg, 10)
	flags |= flagAt(t.flags.ack, 11)
	flags |= flagAt(t.flags.psh, 12)
	flags |= flagAt(t.flags.rst, 13)
	flags |= flagAt(t.flags.syn, 14)
	flags |= flagAt(t.flags.fin, 15)
	encoder.PutUint16(pkt[12:], flags)

	encoder.PutUint16(pkt[14:], t.windowSize)
	encoder.PutUint16(pkt[16:], t.checkSum)
	encoder.PutUint16(pkt[18:], t.urgentPtr)

	return pkt
}
