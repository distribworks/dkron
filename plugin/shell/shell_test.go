package shell

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
