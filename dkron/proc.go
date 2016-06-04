package dkron

import (
	"math/rand"
	"os/exec"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
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
func (a *AgentCommand) invokeJob(execution *Execution) error {
	job := execution.Job

	output, _ := circbuf.NewBuffer(maxBufSize)

	cmd := buildCmd(job)
	cmd.Stderr = output
	cmd.Stdout = output

	// Start a timer to warn about slow handlers
	slowTimer := time.AfterFunc(2*time.Hour, func() {
		log.Warnf("proc: Script '%s' slow, execution exceeding %v", job.Command, 2*time.Hour)
	})

	if err := cmd.Start(); err != nil {
		return err
	}

	// Warn if buffer is overritten
	if output.TotalWritten() > output.Size() {
		log.Warnf("proc: Script '%s' generated %d bytes of output, truncated to %d", job.Command, output.TotalWritten(), output.Size())
	}

	var success bool
	err := cmd.Wait()
	slowTimer.Stop()
	log.WithFields(logrus.Fields{
		"output": output,
	}).Debug("proc: Command output")
	if err != nil {
		log.WithError(err).Error("proc: command error output")
		success = false
	} else {
		success = true
	}

	execution.FinishedAt = time.Now()
	execution.Success = success
	execution.Output = output.Bytes()

	rpcServer, err := a.queryRPCConfig()
	if err != nil {
		return err
	}

	rc := &RPCClient{ServerAddr: string(rpcServer)}
	return rc.callExecutionDone(execution)
}

func (a *AgentCommand) selectServer() serf.Member {
	servers := a.listServers()
	server := servers[rand.Intn(len(servers))]

	return server
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
			log.WithError(err).Fatal("proc: Error parsing command arguments")
		}
		cmd = exec.Command(args[0], args[1:]...)
	}

	return
}
