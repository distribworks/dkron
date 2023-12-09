package main

import (
	"errors"
	"os"
	"os/exec"
	"runtime"

	"github.com/armon/circbuf"
	dkplugin "github.com/distribworks/dkron/v3/plugin"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
	"github.com/mattn/go-shellwords"
)

const (
	windows = "windows"

	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

// reportingWriter This is a Writer implementation that writes back to the host
type reportingWriter struct {
	buffer  *circbuf.Buffer
	cb      dkplugin.StatusHelper
	isError bool
}

func (p reportingWriter) Write(data []byte) (n int, err error) {
	p.cb.Update(data, p.isError)
	return p.buffer.Write(data)
}

// Shell plugin runs shell commands when Execute method is called.
type Shell struct{}

// Execute method of the plugin
func (s *Shell) Execute(args *dktypes.ExecuteRequest, cb dkplugin.StatusHelper) (*dktypes.ExecuteResponse, error) {
	out, err := s.ExecuteImpl(args, cb)
	resp := &dktypes.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// Determine the shell invocation based on OS
func buildCmd(command string, useShell bool, env []string, cwd string) (cmd *exec.Cmd, err error) {
	var shell, flag string

	if useShell {
		if runtime.GOOS == windows {
			shell = "cmd"
			flag = "/C"
		} else {
			shell = "/bin/sh"
			flag = "-c"
		}
		cmd = exec.Command(shell, flag, command)
	} else {
		args, err := shellwords.Parse(command)
		if err != nil {
			return nil, err
		}
		if len(args) == 0 {
			return nil, errors.New("shell: Command missing")
		}
		cmd = exec.Command(args[0], args[1:]...)
	}
	if env != nil {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Dir = cwd
	return
}
