package dkron

import (
	"testing"
)

func TestJob_Validate_MemoryLimit(t *testing.T) {
	tests := []struct {
		name        string
		memLimit    string
		expectError bool
	}{
		{
			name:        "valid memory limit in bytes",
			memLimit:    "1048576",
			expectError: false,
		},
		{
			name:        "valid memory limit with unit",
			memLimit:    "1024MB",
			expectError: false,
		},
		{
			name:        "valid memory limit with lowercase unit",
			memLimit:    "512mb",
			expectError: false,
		},
		{
			name:        "valid memory limit with KB",
			memLimit:    "1024KB",
			expectError: false,
		},
		{
			name:        "valid memory limit with GB", 
			memLimit:    "2GB",
			expectError: false,
		},
		{
			name:        "invalid memory limit format",
			memLimit:    "invalid",
			expectError: true,
		},
		{
			name:        "negative memory limit",
			memLimit:    "-100MB",
			expectError: true,
		},
		{
			name:        "zero memory limit",
			memLimit:    "0",
			expectError: true,
		},
		{
			name:        "empty memory limit",
			memLimit:    "",
			expectError: false, // Empty should be allowed (no limit)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				Name:     "test-memory-job",
				Schedule: "@every 1m",
				Executor: "shell",
				ExecutorConfig: map[string]string{
					"command":   "echo test",
					"mem_limit": tt.memLimit,
				},
			}

			err := job.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for memory limit %q, but got none", tt.memLimit)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for memory limit %q, but got: %v", tt.memLimit, err)
				}
			}
		})
	}
}