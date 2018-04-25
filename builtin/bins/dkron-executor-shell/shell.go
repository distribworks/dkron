package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/armon/circbuf"
	"github.com/mattn/go-shellwords"
	"github.com/victorcoder/dkron/dkron"
)

const (
	windows = "windows"

	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

// FilesOutput plugin that saves each execution log
// in it's own file in the file system.
type Shell struct {
	Param1 string
	Param2 bool
}

// Process method of the plugin
func (s *Shell) Execute(args *dkron.ExecuteRequest) ([]byte, error) {
	output, _ := circbuf.NewBuffer(maxBufSize)

	shell, err := strconv.ParseBool(args.Config["shell"])
	if err != nil {
		shell = false
	}
	command := args.Config["command"]
	env := strings.Split(args.Config["env"], ",")

	cmd, err := buildCmd(command, shell, env)
	if err != nil {
		return nil, err
	}
	cmd.Stderr = output
	cmd.Stdout = output

	// Start a timer to warn about slow handlers
	slowTimer := time.AfterFunc(2*time.Hour, func() {
		log.Printf("shell: Script '%s' slow, execution exceeding %v", command, 2*time.Hour)
	})
	defer slowTimer.Stop()

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
	if err != nil {
		return nil, err
	}

	log.Printf("shell: Command output %s", output)

	return output.Bytes(), nil
}

// Determine the shell invocation based on OS
func buildCmd(command string, useShell bool, env []string) (cmd *exec.Cmd, err error) {
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

	return
}
