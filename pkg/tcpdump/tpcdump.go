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

package tcpdump

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/bowei/lighthouse/pkg/flags"
	"github.com/golang/glog"
)

var (
	// ErrBadSyntax is returned when the tcpdump filter syntax is invalid.
	ErrBadSyntax = errors.New("invalid filter syntax")
)

// Options for tcpdump.
type Options struct {
	Count          int    // -c
	FileSize       int    // -C
	RotateSeconds  int    // -G
	Interface      string // -i
	SnapLen        int    // -s
	OutputFile     string // -w
	FileCountLimit int    // -W
}

// Runner manages the execution of tcpdump.
type Runner struct {
}

// CheckFilter checks the syntax of a BPF filter.
func (r *Runner) CheckFilter(filter string) error {
	cmd := exec.Cmd{
		Args: []string{flags.TCPDumpExecutable, "-d", filter},
		Path: flags.TCPDumpExecutable,
	}
	glog.V(4).Infof("tcpdump = %+v", cmd)
	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.ExitError); ok && !err.Success() {
			return ErrBadSyntax
		}
		return err
	}
	return nil
}

// Run tcpdump.
func (r *Runner) Run(opt *Options, filter string) error {
	if err := r.CheckFilter(filter); err != nil {
		return err
	}

	cmd := exec.Cmd{
		Args: []string{flags.TCPDumpExecutable},
		Path: flags.TCPDumpExecutable,
	}

	if opt.Count != 0 {
		cmd.Args = append(cmd.Args, "-c", fmt.Sprintf("%d", opt.Count))
	}
	if opt.FileSize != 0 {
		cmd.Args = append(cmd.Args, "-C", fmt.Sprintf("%d", opt.FileSize))
	}
	if opt.RotateSeconds != 0 {
		cmd.Args = append(cmd.Args, "-G", fmt.Sprintf("%d", opt.RotateSeconds))
	}
	if opt.Interface != "" {
		cmd.Args = append(cmd.Args, "-i", opt.Interface)
	}
	if opt.SnapLen != 0 {
		cmd.Args = append(cmd.Args, "-s", fmt.Sprintf("%d", opt.SnapLen))
	}
	if opt.OutputFile != "" {
		cmd.Args = append(cmd.Args, "-w", opt.OutputFile)
	}
	if opt.FileCountLimit != 0 {
		cmd.Args = append(cmd.Args, "-W", fmt.Sprintf("%d", opt.FileCountLimit))
	}
	if filter != "" {
		cmd.Args = append(cmd.Args, filter)
	}

	glog.V(4).Infof("tcpdump = %+v", cmd)

	return cmd.Run()
}
