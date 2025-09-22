package shell

import (
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
	dkplugin "github.com/distribworks/dkron/v4/plugin"
	dktypes "github.com/distribworks/dkron/v4/types"
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
			err := processKill(cmd)
			if err != nil {
				jobTimeoutMessage = fmt.Sprintf("shell: Job '%s' execution time exceeding defined timeout %v. SIGKILL returned error. Job may not have been killed", command, jt)
			} else {
				jobTimeoutMessage = fmt.Sprintf("shell: Job '%s' execution time exceeding defined timeout %v. Job was killed", command, jt)
			}

			jobTimedOut = true
		})

		defer slowTimer.Stop()
	}

	// Parse memory limit if specified
	memLimit, err := parseMemoryLimit(args.Config["mem_limit"])
	if err != nil {
		return nil, fmt.Errorf("shell: Error parsing job memory limit: %v", err)
	}

	var memLimitExceededMessage string
	var memLimitExceeded bool

	quit := make(chan int)

	// Start memory monitoring if limit is set
	if memLimit > 0 {
		go func() {
			ticker := time.NewTicker(1 * time.Second) // Check every second
			defer ticker.Stop()
			
			for {
				select {
				case <-quit:
					return
				case <-ticker.C:
					// Get memory usage for the process and its children
					_, totalMem, err := GetTotalCPUMemUsage(cmd.Process.Pid)
					if err != nil {
						// Process might have already finished
						continue
					}
					
					if totalMem > float64(memLimit) {
						// Memory limit exceeded, kill the process
						err := processKill(cmd)
						if err != nil {
							memLimitExceededMessage = fmt.Sprintf("shell: Job '%s' memory usage (%.0f bytes) exceeding defined limit (%d bytes). SIGKILL returned error. Job may not have been killed", command, totalMem, memLimit)
						} else {
							memLimitExceededMessage = fmt.Sprintf("shell: Job '%s' memory usage (%.0f bytes) exceeding defined limit (%d bytes). Job was killed", command, totalMem, memLimit)
						}
						memLimitExceeded = true
						return
					}
				}
			}
		}()
	}

	// FIXME: Debug metrics collection
	// go CollectProcessMetrics(args.JobName, cmd.Process.Pid, quit)

	err = cmd.Wait()
	quit <- cmd.ProcessState.ExitCode()
	close(quit) // exit metric refresh goroutine after job is finished

	if jobTimedOut {
		_, err := output.Write([]byte(jobTimeoutMessage))
		if err != nil {
			log.Printf("Error writing output on timeout event: %v", err)
		}
	}

	if memLimitExceeded {
		_, err := output.Write([]byte(memLimitExceededMessage))
		if err != nil {
			log.Printf("Error writing output on memory limit exceeded event: %v", err)
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

// parseMemoryLimit converts a memory limit string to bytes.
// Accepts formats like "1024", "1024MB", "1GB", "512KB", etc.
// Returns 0 if no limit is specified (empty string).
func parseMemoryLimit(limit string) (int64, error) {
	if limit == "" {
		return 0, nil // No limit
	}

	// Try to parse as a plain number (bytes)
	if value, err := strconv.ParseInt(limit, 10, 64); err == nil {
		if value <= 0 {
			return 0, fmt.Errorf("memory limit must be greater than 0")
		}
		return value, nil
	}

	// Try to parse with units
	limit = strings.ToUpper(strings.TrimSpace(limit))
	
	// Extract the numeric part and unit
	var numStr string
	var unit string
	
	// Find where the number ends and unit begins
	i := 0
	for i < len(limit) && (limit[i] >= '0' && limit[i] <= '9' || limit[i] == '.') {
		i++
	}
	
	if i == 0 {
		return 0, fmt.Errorf("invalid memory limit format: %s", limit)
	}
	
	numStr = limit[:i]
	unit = limit[i:]
	
	// Parse the numeric part
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value in memory limit: %s", numStr)
	}
	
	if value <= 0 {
		return 0, fmt.Errorf("memory limit must be greater than 0")
	}
	
	// Validate and convert unit to bytes
	var multiplier int64
	switch unit {
	case "", "B", "BYTES":
		multiplier = 1
	case "KB", "K":
		multiplier = 1024
	case "MB", "M":
		multiplier = 1024 * 1024
	case "GB", "G":
		multiplier = 1024 * 1024 * 1024
	case "TB", "T":
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unsupported memory unit: %s (supported: B, KB, MB, GB, TB)", unit)
	}
	
	// Check for overflow
	bytes := int64(value * float64(multiplier))
	if bytes <= 0 {
		return 0, fmt.Errorf("memory limit too large or causes overflow")
	}
	
	return bytes, nil
}
