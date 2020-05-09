package main

import (
	"encoding/base64"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

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

// ExecuteImpl do execute command
func (s *Shell) ExecuteImpl(args *dktypes.ExecuteRequest, cb dkplugin.StatusHelper) ([]byte, error) {
	output, _ := circbuf.NewBuffer(maxBufSize)

	shell, err := strconv.ParseBool(args.Config["shell"])
	if err != nil {
		shell = false
	}
	command := args.Config["command"]
	env := strings.Split(args.Config["env"], ",")
	cwd := args.Config["cwd"]

	cmd, err := buildCmd(command, shell, env, cwd)
	if err != nil {
		return nil, err
	}
	err = setCmdAttr(cmd, args.Config)
	if err != nil {
		return nil, err
	}
	// use same buffer for both channels, for the full return at the end
	cmd.Stderr = reportingWriter{buffer: output, cb: cb, isError: true}
	cmd.Stdout = reportingWriter{buffer: output, cb: cb}

	// Start a timer to warn about slow handlers
	slowTimer := time.AfterFunc(2*time.Hour, func() {
		log.Printf("shell: Script '%s' slow, execution exceeding %v", command, 2*time.Hour)
	})
	defer slowTimer.Stop()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	defer stdin.Close()

	payload, err := base64.StdEncoding.DecodeString(args.Config["payload"])
	if err != nil {
		return nil, err
	}

	stdin.Write(payload)
	stdin.Close()

	log.Printf("shell: going to run %s", command)
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// Warn if buffer is overritten
	if output.TotalWritten() > output.Size() {
		log.Printf("shell: Script '%s' generated %d bytes of output, truncated to %d", command, output.TotalWritten(), output.Size())
	}

	err = cmd.Wait()

	// Always log output
	log.Printf("shell: Command output %s", output)

	return output.Bytes(), err
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
		cmd = exec.Command(args[0], args[1:]...)
	}
	if env != nil {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Dir = cwd
	return
}
