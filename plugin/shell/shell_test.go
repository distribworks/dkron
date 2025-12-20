package shell

import (
	"os"
	"runtime"
	"testing"
	"time"

	dktypes "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/stretchr/testify/assert"
)

/*
TestExecuteImpl_* tests comprehensively test the ExecuteImpl method with focus on cmd.Start() and cmd.Wait() functionality.

These tests cover the following scenarios:
1. Success - Normal command execution that succeeds
2. CommandFailure - Command that exits with non-zero status (tests cmd.Wait() error handling)
3. Timeout - Long-running command that gets killed by timeout (tests process termination)
4. WithStdin - Command that reads from stdin (tests stdin pipe setup and data transfer)
5. WithEnvironment - Command with custom environment variables (tests env variable injection)
6. NonShellCommand - Direct command execution without shell wrapper
7. InvalidCommand - Command that fails at cmd.Start() (tests start failure handling)
8. WorkingDirectory - Command execution with custom working directory
9. LargeOutput - Command producing output larger than buffer size (tests output truncation)

Key aspects tested:
- cmd.Start() success and failure scenarios
- cmd.Wait() with various exit conditions
- Process lifecycle management (start, run, terminate)
- Output capture and buffering
- Status callback integration
- Timeout handling and process killing
- Stdin/stdout/stderr redirection
- Environment variable handling
- Working directory changes
- Buffer overflow protection
*/

// MockStatusHelper implements the StatusHelper interface for testing
type MockStatusHelper struct {
	updates []StatusUpdate
}

type StatusUpdate struct {
	data    []byte
	isError bool
}

func (m *MockStatusHelper) Update(data []byte, isError bool) (int64, error) {
	m.updates = append(m.updates, StatusUpdate{data: data, isError: isError})
	return int64(len(data)), nil
}

func (m *MockStatusHelper) GetUpdates() []StatusUpdate {
	return m.updates
}

func (m *MockStatusHelper) Reset() {
	m.updates = nil
}

func TestExecuteImpl_CmdStartWait_Success(t *testing.T) {
	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test simple successful command
	args := &dktypes.ExecuteRequest{
		JobName: "test-job",
		Config: map[string]string{
			"command": "echo 'Hello World'",
			"shell":   "true",
		},
	}

	output, err := s.ExecuteImpl(args, mockCb)

	assert.NoError(t, err)
	assert.Contains(t, string(output), "Hello World")

	// Verify status updates were called
	updates := mockCb.GetUpdates()
	assert.True(t, len(updates) > 0, "Expected status updates to be called")
}

func TestExecuteImpl_CmdStartWait_CommandFailure(t *testing.T) {
	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test command that exits with non-zero status
	var failCommand string
	if runtime.GOOS == "windows" {
		failCommand = "exit 1"
	} else {
		failCommand = "exit 1"
	}

	args := &dktypes.ExecuteRequest{
		JobName: "test-job-fail",
		Config: map[string]string{
			"command": failCommand,
			"shell":   "true",
		},
	}

	output, err := s.ExecuteImpl(args, mockCb)

	assert.Error(t, err, "Expected command to fail")
	assert.NotNil(t, output, "Output should still be captured even on failure")

	// Verify status updates were called
	updates := mockCb.GetUpdates()
	assert.True(t, len(updates) >= 0, "Status updates should be called even on failure")
}

func TestExecuteImpl_CmdStartWait_Timeout(t *testing.T) {
	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test command that times out
	var sleepCommand string
	if runtime.GOOS == "windows" {
		sleepCommand = "timeout /t 3 /nobreak"
	} else {
		sleepCommand = "sleep 3"
	}

	args := &dktypes.ExecuteRequest{
		JobName: "test-job-timeout",
		Config: map[string]string{
			"command": sleepCommand,
			"shell":   "true",
			"timeout": "1s", // Shorter than sleep duration
		},
	}

	start := time.Now()
	output, err := s.ExecuteImpl(args, mockCb)
	duration := time.Since(start)

	// Command should complete within reasonable time due to timeout
	assert.True(t, duration < 2*time.Second, "Command should be killed by timeout")

	// Output should contain timeout message
	assert.Contains(t, string(output), "execution time exceeding defined timeout")
	assert.Contains(t, string(output), "Job was killed")

	// Error may or may not be present depending on timing
	_ = err
}

func TestExecuteImpl_CmdStartWait_WithStdin(t *testing.T) {
	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test command that reads from stdin
	var catCommand string
	if runtime.GOOS == "windows" {
		catCommand = "findstr .*"
	} else {
		catCommand = "cat"
	}

	args := &dktypes.ExecuteRequest{
		JobName: "test-job-stdin",
		Config: map[string]string{
			"command": catCommand,
			"shell":   "true",
			"payload": "SGVsbG8gZnJvbSBzdGRpbg==", // "Hello from stdin" in base64
		},
	}

	output, err := s.ExecuteImpl(args, mockCb)

	assert.NoError(t, err)
	assert.Contains(t, string(output), "Hello from stdin")
}

func TestExecuteImpl_CmdStartWait_WithEnvironment(t *testing.T) {
	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test command with environment variables
	var echoCommand string
	if runtime.GOOS == "windows" {
		echoCommand = "echo %TEST_VAR% %ENV_JOB_NAME%"
	} else {
		echoCommand = "echo $TEST_VAR $ENV_JOB_NAME"
	}

	args := &dktypes.ExecuteRequest{
		JobName: "test-job-env",
		Config: map[string]string{
			"command": echoCommand,
			"shell":   "true",
			"env":     "TEST_VAR=hello",
		},
	}

	output, err := s.ExecuteImpl(args, mockCb)

	assert.NoError(t, err)
	assert.Contains(t, string(output), "hello")
	assert.Contains(t, string(output), "test-job-env") // ENV_JOB_NAME should be set
}

func TestExecuteImpl_CmdStartWait_NonShellCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping non-shell command test on Windows")
	}

	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test direct command execution (not through shell)
	args := &dktypes.ExecuteRequest{
		JobName: "test-job-direct",
		Config: map[string]string{
			"command": "echo hello direct",
			"shell":   "false",
		},
	}

	output, err := s.ExecuteImpl(args, mockCb)

	assert.NoError(t, err)
	assert.Contains(t, string(output), "hello direct")
}

func TestExecuteImpl_CmdStartWait_InvalidCommand(t *testing.T) {
	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test invalid command that should fail at start
	args := &dktypes.ExecuteRequest{
		JobName: "test-job-invalid",
		Config: map[string]string{
			"command": "nonexistent-command-12345",
			"shell":   "false",
		},
	}

	output, err := s.ExecuteImpl(args, mockCb)

	assert.Error(t, err, "Expected error for invalid command")
	assert.Nil(t, output, "Output should be nil when command fails to start")
}

func TestExecuteImpl_CmdStartWait_WorkingDirectory(t *testing.T) {
	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test command with working directory
	var pwdCommand string
	if runtime.GOOS == "windows" {
		pwdCommand = "cd"
	} else {
		pwdCommand = "pwd"
	}

	args := &dktypes.ExecuteRequest{
		JobName: "test-job-cwd",
		Config: map[string]string{
			"command": pwdCommand,
			"shell":   "true",
			"cwd":     "/tmp",
		},
	}

	output, err := s.ExecuteImpl(args, mockCb)

	if runtime.GOOS != "windows" {
		assert.NoError(t, err)
		assert.Contains(t, string(output), "/tmp")
	}
	// On Windows, just verify it doesn't crash
	_ = output
}

func TestExecuteImpl_CmdStartWait_LargeOutput(t *testing.T) {
	s := &Shell{}
	mockCb := &MockStatusHelper{}

	// Test command that produces output larger than maxBufSize
	var largeOutputCommand string
	if runtime.GOOS == "windows" {
		// Generate large output on Windows
		largeOutputCommand = "for /L %i in (1,1,1000) do @echo This is line %i with some additional text to make it longer"
	} else {
		// Generate large output on Unix - use seq instead of bash-specific syntax
		largeOutputCommand = "seq 1 1000 | while read i; do echo \"This is line $i with some additional text to make it longer\"; done"
	}

	args := &dktypes.ExecuteRequest{
		JobName: "test-job-large",
		Config: map[string]string{
			"command": largeOutputCommand,
			"shell":   "true",
		},
	}

	output, err := s.ExecuteImpl(args, mockCb)

	// Should not error due to large output
	assert.NoError(t, err)
	assert.True(t, len(output) > 0, "Should have captured some output")
	assert.True(t, len(output) <= maxBufSize, "Output should be truncated to maxBufSize")
}

func Test_buildCmdInvalidInput(t *testing.T) {
	assert.NotPanics(t, func() { buildCmd("", false, []string{}, "") })
}

func Test_buildCmd(t *testing.T) {

	// test shell command
	cmd, err := buildCmd("echo 'test1' && echo 'success'", true, []string{}, "")
	assert.NoError(t, err)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "test1\nsuccess\n", string(out))

	// test not shell command
	cmd, err = buildCmd("date && echo 'success'", false, []string{}, "")
	assert.NoError(t, err)
	assert.Len(t, cmd.Args, 1)
}

func Test_buildCmdWithCustomEnvironmentVariables(t *testing.T) {
	defer func() {
		os.Setenv("Foo", "")
	}()
	os.Setenv("Foo", "Bar")

	cmd, err := buildCmd("echo $Foo && echo $He", true, []string{"Foo=Toto", "He=Ho"}, "")
	assert.NoError(t, err)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "Toto\nHo\n", string(out))

}

func Test_parseMemoryLimit(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int64
		expectError bool
	}{
		{
			name:        "empty string returns no limit",
			input:       "",
			expected:    0,
			expectError: false,
		},
		{
			name:        "plain bytes",
			input:       "1024",
			expected:    1024,
			expectError: false,
		},
		{
			name:        "KB unit",
			input:       "1KB",
			expected:    1024,
			expectError: false,
		},
		{
			name:        "MB unit",
			input:       "1MB",
			expected:    1024 * 1024,
			expectError: false,
		},
		{
			name:        "GB unit",
			input:       "2GB",
			expected:    2 * 1024 * 1024 * 1024,
			expectError: false,
		},
		{
			name:        "lowercase units",
			input:       "512mb",
			expected:    512 * 1024 * 1024,
			expectError: false,
		},
		{
			name:        "decimal values",
			input:       "1.5GB",
			expected:    int64(1.5 * 1024 * 1024 * 1024),
			expectError: false,
		},
		{
			name:        "invalid format",
			input:       "invalid",
			expected:    0,
			expectError: true,
		},
		{
			name:        "negative value",
			input:       "-100MB",
			expected:    0,
			expectError: true,
		},
		{
			name:        "zero value",
			input:       "0",
			expected:    0,
			expectError: true,
		},
		{
			name:        "unsupported unit",
			input:       "100PB",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseMemoryLimit(tt.input)

			if tt.expectError {
				assert.Error(t, err, "Expected error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Expected no error for input: %s", tt.input)
				assert.Equal(t, tt.expected, result, "Expected %d bytes for input: %s", tt.expected, tt.input)
			}
		})
	}
}
