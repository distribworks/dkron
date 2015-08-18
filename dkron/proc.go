package dkron

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/armon/circbuf"
	"github.com/hashicorp/serf/serf"
)

const (
	windows = "windows"

	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 8 * 1024
)

// spawn command that specified as proc.
func spawnProc(proc string) (*exec.Cmd, error) {
	cs := []string{"/bin/bash", "-c", proc}
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ())

	log.Printf("Starting %s\n", proc)
	err := cmd.Start()
	if err != nil {
		log.Errorf("Failed to start %s: %s\n", proc, err)
		return nil, err
	}
	return cmd, nil
}

// invokeJob will execute the given job. Depending on the event.
func (a *AgentCommand) invokeJob(job *Job, execution *Execution) error {
	output, _ := circbuf.NewBuffer(maxBufSize)

	// Determine the shell invocation based on OS
	var shell, flag string
	if runtime.GOOS == windows {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

	cmd := exec.Command(shell, flag, job.Command)
	cmd.Stderr = output
	cmd.Stdout = output

	// Start a timer to warn about slow handlers
	slowTimer := time.AfterFunc(2*time.Hour, func() {
		log.Warnf("Script '%s' slow, execution exceeding %v", job.Command, 2*time.Hour)
	})

	if err := cmd.Start(); err != nil {
		return err
	}

	// Warn if buffer is overritten
	if output.TotalWritten() > output.Size() {
		log.Warnf("Script '%s' generated %d bytes of output, truncated to %d", job.Command, output.TotalWritten(), output.Size())
	}

	var success bool
	err := cmd.Wait()
	slowTimer.Stop()
	log.Debugf("Command output: %s", output)
	if err != nil {
		log.Error(err)
		success = false
	} else {
		success = true
	}

	execution.FinishedAt = time.Now()
	execution.Success = success
	execution.Output = output.Bytes()

	executionJson, _ := json.Marshal(execution)

	params := &serf.QueryParam{
		FilterTags: map[string]string{"server": "true"},
		RequestAck: true,
	}

	qr, err := a.serf.Query(QueryExecutionDone, executionJson, params)
	if err != nil {
		log.WithFields(logrus.Fields{
			"query": QueryExecutionDone,
			"error": err,
		}).Debug("Error sending query")
	}
	defer qr.Close()

	ackCh := qr.AckCh()
	respCh := qr.ResponseCh()

	for !qr.Finished() {
		select {
		case ack, ok := <-ackCh:
			if ok {
				log.WithFields(logrus.Fields{
					"query": QueryExecutionDone,
					"from":  ack,
				}).Debug("Received ack")
			}
		case resp, ok := <-respCh:
			if ok {
				log.WithFields(logrus.Fields{
					"query":   QueryExecutionDone,
					"from":    resp.From,
					"payload": string(resp.Payload),
				}).Debug("Received response")
			}
		}
	}
	log.Debugf("Done receiving acks and responses from %s query", QueryExecutionDone)

	return nil
}
