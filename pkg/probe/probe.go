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
	"net"

	"github.com/golang/glog"
)

// SendTCP TODO
func SendTCP(src string, srcPort int, dest string, destPort int, magic string) error {
	srcAddr, err := net.ResolveIPAddr("ip4", src)
	if err != nil {
		return err
	}
	destAddr, err := net.ResolveIPAddr("ip4", dest)
	if err != nil {
		return err
	}
	conn, err := net.DialIP("ip4:tcp", srcAddr, destAddr)
	if err != nil {
		return err
	}

	tcp := &tcpPacket{
		srcPort:  uint16(srcPort),
		destPort: uint16(destPort),
		flags:    tcpFlags{syn: true},
		seq:      1,
	}

	pkt := make([]byte, 1024)
	len := tcp.encode(pkt, net.ParseIP(src), net.ParseIP(dest), []byte{})
	glog.V(2).Infof("Encoded TCP (%d bytes): %v", len, pkt[:len])
	len, err = conn.Write(pkt[:len])
	glog.V(2).Infof("conn.Write(pkt) = %d, %v", len, err)

	return err
}
