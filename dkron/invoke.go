package dkron

import (
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/armon/circbuf"
	"github.com/hashicorp/serf/serf"
	"github.com/mattn/go-shellwords"
)

const (
	windows = "windows"

	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

// invokeJob will execute the given job. Depending on the event.
func (a *Agent) invokeJob(job *Job, execution *Execution) error {
	output, _ := circbuf.NewBuffer(maxBufSize)

	var success bool

	jex := job.Executor
	exc := job.ExecutorConfig
	if jex == "" {
		jex = "shell"
		exc = map[string]string{
			"command": job.Command,
			"shell":   strconv.FormatBool(job.Shell),
			"env":     strings.Join(job.EnvironmentVariables, ","),
		}
		log.Warning("invoke: Deprecation waring! fields command, " +
			"shell and environment_variables params are deprecated and will be removed in future versions. " +
			"Consider migrating the job definition to the shell executor plugin")
	}

	// Check if executor is exists
	if executor, ok := a.ExecutorPlugins[jex]; ok {
		log.WithField("plugin", jex).Debug("invoke: calling executor plugin")
		runningExecutions.Store(execution.GetGroup(), execution)
		out, err := executor.Execute(&ExecuteRequest{
			JobName: job.Name,
			Config:  exc,
		})

		if err == nil && out.Error != "" {
			err = errors.New(out.Error)
		}
		if err != nil {
			log.WithError(err).WithField("job", job.Name).WithField("plugin", executor).Error("invoke: command error output")
			success = false
			output.Write([]byte(err.Error() + "\n"))
		} else {
			success = true
		}

		if out != nil {
			output.Write(out.Output)
		}
	} else {
		log.WithField("executor", jex).Error("invoke: Specified executor is not present")
	}

	execution.FinishedAt = time.Now()
	execution.Success = success
	execution.Output = output.Bytes()

	rpcServer, err := a.getServerRPCAddresFromTags()
	if err != nil {
		return err
	}

	runningExecutions.Delete(execution.GetGroup())

	return a.GRPCClient.CallExecutionDone(rpcServer, execution)
}

func (a *Agent) selectServer() serf.Member {
	servers := a.listServers()
	server := servers[rand.Intn(len(servers))]

	return server
}

func (a *Agent) getServerRPCAddresFromTags() (string, error) {
	s := a.selectServer()

	if addr, ok := s.Tags["dkron_rpc_addr"]; ok {
		return addr, nil
	}
	return "", ErrNoRPCAddress
}

// Determine the shell invocation based on OS
func buildCmd(job *Job) (cmd *exec.Cmd) {
	var shell, flag string

	if job.Shell {
		if runtime.GOOS == windows {
			shell = "cmd"
			flag = "/C"
		} else {
			shell = "/bin/sh"
			flag = "-c"
		}
		cmd = exec.Command(shell, flag, job.Command)
	} else {
		args, err := shellwords.Parse(job.Command)
		if err != nil {
			log.WithError(err).Fatal("invoke: Error parsing command arguments")
		}
		cmd = exec.Command(args[0], args[1:]...)
	}
	if job.EnvironmentVariables != nil {
		cmd.Env = append(os.Environ(), job.EnvironmentVariables...)
	}

	return
}
