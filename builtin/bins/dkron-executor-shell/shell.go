package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
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
	"github.com/struCoder/pidusage"
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

	jobTimeout := args.Config["timeout"]
	var jt time.Duration

	if jobTimeout != "" {
		jt, err = time.ParseDuration(jobTimeout)
		if err != nil {
			return nil, errors.New("shell: Error parsing job timeout")
		}
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	var jobTimeoutMessage string
	var jobTimedOut bool

	if jt != 0 {
		slowTimer := time.AfterFunc(jt, func() {
			err = cmd.Process.Kill()
			if err != nil {
				jobTimeoutMessage = fmt.Sprintf("shell: Job '%s' execution time exceeding defined timeout %v. SIGKILL returned error. Job may not have been killed", command, jt)
			} else {
				jobTimeoutMessage = fmt.Sprintf("shell: Job '%s' execution time exceeding defined timeout %v. Job was killed", command, jt)
			}

			jobTimedOut = true
			return
		})

		defer slowTimer.Stop()
	}

	// Warn if buffer is overwritten
	if output.TotalWritten() > output.Size() {
		log.Printf("shell: Script '%s' generated %d bytes of output, truncated to %d", command, output.TotalWritten(), output.Size())
	}

	pid := cmd.Process.Pid
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				stat, err := pidusage.GetStat(pid)
				if err != nil {
					log.Printf("Error getting pid statistics: %v", err)
					return
				}

				mem, err := calculateMemory(pid)
				if err != nil {
					log.Printf("Error calculating memory metrics: %v", err)
					return
				}

				cpu := stat.CPU
				updateMetric(args.JobName, memUsage, float64(mem))
				updateMetric(args.JobName, cpuUsage, cpu)
				time.Sleep(1 * time.Second) // Refreshing metrics in real-time each second
			}
		}
	}()

	err = cmd.Wait()
	close(quit) // exit metric refresh goroutine after job is finished

	if jobTimedOut {
		_, err := output.Write([]byte(jobTimeoutMessage))
		if err != nil {
			log.Printf("Error writing output on timeout event: %v", err)
		}
	}

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

func calculateMemory(pid int) (uint64, error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/smaps", pid))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	res := uint64(0)
	rfx := []byte("Rss:")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, rfx) {
			var size uint64
			_, err := fmt.Sscanf(string(line[4:]), "%d", &size)
			if err != nil {
				return 0, err
			}
			res += size
		}
	}
	if err := r.Err(); err != nil {
		return 0, err
	}
	return res, nil
}
