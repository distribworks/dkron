//go:build windows
// +build windows

package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/armon/circbuf"
	dkplugin "github.com/distribworks/dkron/v3/plugin"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
)

func setCmdAttr(cmd *exec.Cmd, config map[string]string) error {
	return nil
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

	executionInfo := strings.Split(fmt.Sprintf("ENV_JOB_NAME=%s", args.JobName), ",")
	env = append(env, executionInfo...)

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

	jobTimeout := args.Config["timeout"]
	var jt time.Duration

	if jobTimeout != "" {
		jt, err = time.ParseDuration(jobTimeout)
		if err != nil {
			return nil, errors.New("shell: Error parsing job timeout")
		}
	}

	log.Printf("shell: going to run %s", command)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	var jobTimeoutMessage string
	var jobTimedOut bool

	if jt != 0 {
		slowTimer := time.AfterFunc(jt, func() {
			// Kill child process to avoid cmd.Wait()
			err = cmd.Process.Kill()
			if err != nil {
				jobTimeoutMessage = fmt.Sprintf("shell: Job '%s' execution time exceeding defined timeout %v. SIGKILL returned error. Job may not have been killed", command, jt)
			} else {
				jobTimeoutMessage = fmt.Sprintf("shell: Job '%s' execution time exceeding defined timeout %v. Job was killed", command, jt)
			}

			jobTimedOut = true
		})

		defer slowTimer.Stop()
	}

	quit := make(chan int)

	go CollectProcessMetrics(args.JobName, cmd.Process.Pid, quit)

	err = cmd.Wait()
	quit <- cmd.ProcessState.ExitCode()
	close(quit) // exit metric refresh goroutine after job is finished

	if jobTimedOut {
		_, err := output.Write([]byte(jobTimeoutMessage))
		if err != nil {
			log.Printf("Error writing output on timeout event: %v", err)
		}
	}

	// Warn if buffer is overwritten
	if output.TotalWritten() > output.Size() {
		log.Printf("shell: Script '%s' generated %d bytes of output, truncated to %d", command, output.TotalWritten(), output.Size())
	}

	// Always log output
	log.Printf("shell: Command output %s", output)

	return output.Bytes(), err
}
