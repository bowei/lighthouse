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
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestTCPEncode(t *testing.T) {
	t.Parallel()

	src := net.ParseIP("1.1.1.1")
	dest := net.ParseIP("2.2.2.2")

	for _, tc := range []struct {
		desc string
		tcp  tcpHeader
		data []byte
		want []byte
	}{
		{
			desc: "localhost",
			tcp: tcpHeader{
				srcPort:    80,
				destPort:   81,
				seq:        0x1,
				ack:        0x2,
				windowSize: 0x3,
				flags: tcpFlags{
					syn: true,
				},
			},
			want: []byte{
				0, 80, 0, 81, 0, 0, 0, 1, 0, 0,
				0, 2, 64, 0, 0, 3, 217, 234, 0, 0,
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			pkt := make([]byte, TCPHeaderSize)
			tc.tcp.Encode(pkt, src, dest, tc.data)
			if !reflect.DeepEqual(tc.want, pkt) {
				t.Errorf("tcp.Encode() = %v, want %v; tcp = %+v, data = %v", pkt, tc.want, tc.tcp, tc.data)
			}
		})
	}
}

func TestTCPChecksummer(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		desc string
		data [][]byte
		want uint16
	}{
		{
			desc: "empty",
			data: [][]byte{},
			want: 0xffff,
		},
		{
			desc: "1 byte",
			data: [][]byte{[]byte{0x55}},
			want: 0xffaa,
		},
		{
			desc: "2 bytes",
			data: [][]byte{[]byte{0x55, 0x88}},
			want: 0x77aa,
		},
		{
			desc: "3 bytes",
			data: [][]byte{[]byte{0x55, 0x88, 0x99}},
			want: 0x7711,
		},
		{
			desc: "3 bytes / 1 at a time",
			data: [][]byte{[]byte{0x55}, []byte{0x88}, []byte{0x99}},
			want: 0x7711,
		},
		{
			desc: "3 bytes / 2 1",
			data: [][]byte{[]byte{0x55, 0x88}, []byte{0x99}},
			want: 0x7711,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			fmt.Println("---")
			c := &tcpChecksumer{}
			for _, b := range tc.data {
				c.add(b)
			}
			got := c.finalize()
			if got != tc.want {
				t.Errorf("c.finalize() = %x, want %x; bytes: %v", got, tc.want, tc.data)
			}
		})
	}
}
