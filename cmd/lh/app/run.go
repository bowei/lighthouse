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
	"sort"
)

type subcommand interface {
	flags() *flag.FlagSet
	run() int
}

var allSubcommands = map[string]subcommand{}

// Run the app.
func Run() {
	if len(os.Args) < 2 {
		fmt.Println("Need a subcommand")
		os.Exit(1)
	}

	cmd := os.Args[1]
	var (
		sc subcommand
		ok bool
	)
	if sc, ok = allSubcommands[cmd]; !ok {
		fmt.Printf("Invalid subcommand %q. Available subcommands:\n", cmd)
		var names []string
		for name := range allSubcommands {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Printf("  %s\n", name)
		}
		os.Exit(1)
	}

	f := sc.flags()
	mergeGlobalFlags(f)

	args := os.Args[2:]
	if err := f.Parse(args); err != nil {
		f.Usage()
		os.Exit(1)
	}
	flag.Parse()

	os.Exit(sc.run())
}
