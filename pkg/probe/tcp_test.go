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
	"reflect"
	"testing"
)

func TestTCPEncode(t *testing.T) {
	for _, tc := range []struct {
		desc string
		tcp  tcpHeader
		want []byte
	}{
		{
			desc: "basic",
			tcp: tcpHeader{
				srcPort:  80,
				destPort: 0xff0f,
				seq:      1,
				ack:      2,
				flags: tcpFlags{
					syn: true,
				},
			},
			want: []byte{0, 80, 0xff, 0xf, 0, 0, 0, 1, 0, 0, 0, 2, 64, 0, 0, 0, 0, 0, 0, 0},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			pkt := make([]byte, TCPHeaderSize)
			tc.tcp.Encode(pkt)
			if !reflect.DeepEqual(tc.want, pkt) {
				t.Errorf("tcp.Encode() = %v, want %v; tcp = %+v", pkt, tc.want, tc.tcp)
			}
		})
	}
}
