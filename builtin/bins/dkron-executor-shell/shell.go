package main

import (
	"encoding/base64"
	"fmt"
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
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"
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
	if command == "" {
		return nil, err
	}

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
	jobMemLimit := args.Config["mem_limit_kb"]

	if jobTimeout == "" {
		log.Infof("shell: Job '%v' doesn't have configured timeout. Defaulting to 24h", args.JobName)
		jobTimeout = "24h"
	}

	if args.Config["mem_limit_kb"] == "" {
		log.Infof("shell: Job '%v' doesn't have configured mem_limit_kb", args.JobName)
		jobMemLimit = "inf"
	}

	t, err := time.ParseDuration(jobTimeout)
	if err != nil {
		log.Infof("shell: Job '%v' can't parse job timeout. Defaulting to 24h", args.JobName)
		jobTimeout = "24h"
	}

	startTime := time.Now()
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	slowTimer := time.AfterFunc(t, func() {
		j := fmt.Sprintf("shell: Job '%s' execution time exceeding timeout %v. Killing job.", command, t)
		output.Write([]byte(j))
		cmd.Process.Kill()

	})
	defer slowTimer.Stop()

	// Warn if buffer is overritten
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
				stat, _ := pidusage.GetStat(pid)
				mem, _ := calculateMemory(pid)
				cpu := stat.CPU
				updateMetric(args.JobName, memUsage, float64(mem))
				updateMetric(args.JobName, cpuUsage, cpu)
				if jobMemLimit != "inf" {
					i, _ := strconv.ParseUint(jobMemLimit, 0, 64)
					if mem > i {
						j := fmt.Sprintf("shell: Job '%s' memory limit exceeded %vkb. Killing job.", command, i)
						output.Write([]byte(j))
						cmd.Process.Kill()
						return
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	err = cmd.Wait()
	close(quit)

	executionTime.Set(time.Since(startTime).Seconds())
	exitCode.Set(float64(cmd.ProcessState.ExitCode()))
	lastExecutionTimestamp.Set(float64(time.Now().Unix()))

	push.New(getEnv("PUSHGATEWAY_URL"), "dkron_job_push").Collector(executionTime).Grouping("job_name", args.JobName).Add()
	push.New(getEnv("PUSHGATEWAY_URL"), "dkron_job_push").Collector(exitCode).Grouping("job_name", args.JobName).Add()
	push.New(getEnv("PUSHGATEWAY_URL"), "dkron_job_push").Collector(lastExecutionTimestamp).Grouping("job_name", args.JobName).Add()

	updateMetric(args.JobName, memUsage, 0)
	updateMetric(args.JobName, cpuUsage, 0)

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
		if len(args) == 0 {
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
