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
)

func SendTCP() {
	netaddr, _ := net.ResolveIPAddr("ip4", "127.0.0.1")
	conn, err := net.DialIP("ip4:tcp", nil, netaddr)
	if err != nil {
		panic(err)
	}
	tcp := &tcpHeader{
		srcPort:  80,
		destPort: 8080,
		flags:    tcpFlags{syn: true},
		seq:      1,
	}

	pkt := make([]byte, 1024)
	len := tcp.Encode(pkt, net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"), []byte{})
	fmt.Printf("Encode(...) = %v\n", len)

	fmt.Printf("pkt = %v\n", pkt[:len])
	fmt.Printf("[")
	for _, x := range pkt[:len] {
		fmt.Printf("%02x ", x)
	}
	fmt.Printf("]\n")
	len, err = conn.Write(pkt[:len])
	fmt.Printf("%v, %v\n", len, err)
}
