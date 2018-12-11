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
	"fmt"
	"net"
)

type tcpFlags struct {
	ns, cwr, ece, urg, ack, psh, rst, syn, fin bool
}

const TCPHeaderSize = 20

func newTCPHeader() *tcpHeader {
	return &tcpHeader{dataOffset: 5}
}

type tcpHeader struct {
	srcPort    uint16 // 0
	destPort   uint16 // 2
	seq        uint32 // 4
	ack        uint32 // 8
	dataOffset uint16
	flags      tcpFlags // 12
	windowSize uint16   // 14
	urgentPtr  uint16   // 18
}

func (t *tcpHeader) Encode(pkt []byte, src, dest net.IP, data []byte) {
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
	encoder.PutUint16(pkt[18:], t.urgentPtr)
	// Checksum is last (compute with pseudoheader and zeros).
	encoder.PutUint16(pkt[16:], checksumTCP(src, dest, pkt[:20], data))
}

const tcpProtoNum = 6

func checksumTCP(src, dest net.IP, tcpHeader, data []byte) uint16 {
	chk := &tcpChecksumer{}
	// Encode pseudoheader.
	chk.add(src.To4())
	chk.add(dest.To4())
	var pseudoHeader [4]byte
	enc := binary.BigEndian
	enc.PutUint16(pseudoHeader[:2], tcpProtoNum<<8)
	enc.PutUint16(pseudoHeader[2:], uint16(len(data)))
	chk.add(pseudoHeader[:])

	chk.add(tcpHeader)
	chk.add(data)

	return chk.finalize()
}

type tcpChecksumer struct {
	sum     uint32
	oddByte byte
	length  int
}

func (c *tcpChecksumer) finalize() uint16 {
	ret := c.sum
	if c.length%2 > 0 {
		ret += uint32(c.oddByte)
	}
	for ret>>16 > 0 {
		ret = ret&0xffff + ret>>16
	}
	return ^uint16(ret)
}

func (c *tcpChecksumer) add(data []byte) {
	if len(data) == 0 {
		return
	}
	haveOddByte := c.length%2 > 0
	c.length += len(data)
	if haveOddByte {
		data = append([]byte{c.oddByte}, data...)
	}
	fmt.Printf("DATA %v\n", data)
	for i := 0; i < len(data)-1; i += 2 {
		c.sum += uint32(data[0]) + uint32(data[1])<<8
	}
	if c.length%2 > 0 {
		fmt.Printf("updating oddByte\n")
		c.oddByte = data[len(data)-1]
	}
}
