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
	"fmt"
	"os"

	"github.com/golang/glog"
)

type subcommand interface {
	name() string
	run(args []string) int
}

var allSubcommands []subcommand

// Run the app.
func Run() {
	if flag.NArg() > len(os.Args) {
		fmt.Println("Need a subcommand")
	}
	flag.Parse()
	cmd := os.Args[flag.NArg()]
	rest := os.Args[flag.NArg()+1:]
	glog.V(2).Infof("cmd=%q rest=%v", cmd, rest)
	for _, sc := range allSubcommands {
		if sc.name() == cmd {
			os.Exit(sc.run(rest))
		}
	}

	fmt.Printf("Invalid subcommand %q\n", cmd)
}
