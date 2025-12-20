package dkron

import (
	"context"
	"io"
)

// Storage is the interface that should be used by any
// storage engine implemented for dkron. It contains the
// minimum set of operations that are needed to have a working
// dkron store.
type Storage interface {
	SetJob(ctx context.Context, job *Job, copyDependentJobs bool) error
	DeleteJob(ctx context.Context, name string) (*Job, error)
	DeleteExecutions(ctx context.Context, jobName string) error
	SetExecution(ctx context.Context, execution *Execution) (string, error)
	SetExecutionDone(ctx context.Context, execution *Execution) (bool, error)
	GetJobs(ctx context.Context, options *JobOptions) ([]*Job, error)
	GetJob(ctx context.Context, name string, options *JobOptions) (*Job, error)
	GetExecution(ctx context.Context, jobName string, executionName string) (*Execution, error)
	GetExecutions(ctx context.Context, jobName string, opts *ExecutionOptions) ([]*Execution, error)
	GetRunningExecutions(ctx context.Context, jobName string) ([]*Execution, error)
	GetExecutionGroup(ctx context.Context, execution *Execution, opts *ExecutionOptions) ([]*Execution, error)
	GetGroupedExecutions(ctx context.Context, jobName string, opts *ExecutionOptions) (map[int64][]*Execution, []int64, error)
	Shutdown() error
	Snapshot(w io.WriteCloser) error
	Restore(r io.ReadCloser) error
}
