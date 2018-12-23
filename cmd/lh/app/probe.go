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

package app

import (
	"flag"

	"github.com/bowei/lighthouse/pkg/probe"
	"github.com/golang/glog"
)

var (
	probeFlagSet = flag.NewFlagSet("probe", flag.ExitOnError)
	probeFlags   = struct {
		endpoint *string
		port     *int
		magic    *string
	}{
		endpoint: probeFlagSet.String("endpoint", "", "endpoint to send to"),
		port:     probeFlagSet.Int("port", 80, "port to send to"),
		magic:    probeFlagSet.String("magic", "magic", "magic packet identity"),
	}
)

func init() {
	allSubcommands = append(allSubcommands, &probeCommand{})
}

type probeCommand struct{}

func (c *probeCommand) name() string {
	return "probe"
}

func (c *probeCommand) run(args []string) int {
	probeFlagSet.Parse(args)
	glog.V(2).Infof("runProbe endpoint=%s magic=%s", *probeFlags.endpoint, *probeFlags.magic)
	err := probe.SendTCP("127.0.0.1", 3000, *probeFlags.endpoint, *probeFlags.port, *probeFlags.magic)
	glog.V(2).Infof("probe.SendTCP = %v", err)

	return 0
}
