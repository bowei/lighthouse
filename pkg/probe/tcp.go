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
	"net"
)

const (
	tcpHeaderSize = 20
	tcpProtoNum   = 6
)

type tcpFlags struct {
	ns, cwr, ece, urg, ack, psh, rst, syn, fin bool
}

type tcpPacket struct {
	srcPort    uint16 // 0
	destPort   uint16 // 2
	seq        uint32 // 4
	ack        uint32 // 8
	dataOffset uint16
	flags      tcpFlags // 12
	windowSize uint16   // 14
	urgentPtr  uint16   // 18
}

func (t *tcpPacket) Encode(pkt []byte, src, dest net.IP, data []byte) int {
	encoder := binary.BigEndian
	encoder.PutUint16(pkt, t.srcPort)
	encoder.PutUint16(pkt[2:], t.destPort)
	encoder.PutUint32(pkt[4:], t.seq)
	encoder.PutUint32(pkt[8:], t.ack)

	if t.dataOffset == 0 {
		// If nil-initialized, then assume standard size (5*32 bits).
		pkt[12] = 5 << 4
	} else {
		pkt[12] = uint8(t.dataOffset&0xf) << 4
	}

	if t.flags.ns {
		pkt[12] |= 1 << 7
	}

	var flags uint8
	for offset, f := range []bool{
		t.flags.fin,
		t.flags.syn,
		t.flags.rst,
		t.flags.psh,
		t.flags.ack,
		t.flags.urg,
		t.flags.ece,
		t.flags.cwr,
	} {
		if f {
			flags |= 1 << uint(offset)
		}
	}
	pkt[13] = flags

	encoder.PutUint16(pkt[14:], t.windowSize)
	encoder.PutUint16(pkt[18:], t.urgentPtr)
	// Checksum is last (compute with pseudoheader and zeros).
	checksum := checksumTCP(src, dest, pkt[:tcpHeaderSize], data)
	pkt[16] = uint8(checksum & 0xff)
	pkt[17] = uint8(checksum >> 8)

	copy(pkt[tcpHeaderSize:], data)

	return tcpHeaderSize + len(data)
}

func checksumTCP(src, dest net.IP, tcpHeader, data []byte) uint16 {
	chk := &tcpChecksumer{}

	// Encode pseudoheader.
	chk.add(src.To4())
	chk.add(dest.To4())

	var pseudoHeader [4]byte
	pseudoHeader[1] = tcpProtoNum
	binary.BigEndian.PutUint16(pseudoHeader[2:], uint16(len(data)+len(tcpHeader)))
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
	for i := 0; i < len(data)-1; i += 2 {
		c.sum += uint32(data[i]) + uint32(data[i+1])<<8
	}
	if c.length%2 > 0 {
		c.oddByte = data[len(data)-1]
	}
}
