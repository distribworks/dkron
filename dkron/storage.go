package dkron

import "io"

// Storage is the interface that should be used by any
// storage engine implemented for dkron. It contains the
// minumum set of operations that are needed to have a working
// dkron store.
type Storage interface {
	SetJob(job *Job, copyDependentJobs bool) error
	DeleteJob(name string) (*Job, error)
	SetExecution(execution *Execution) (string, error)
	SetExecutionDone(execution *Execution) (bool, error)
	GetJobs(options *JobOptions) ([]*Job, error)
	GetJob(name string, options *JobOptions) (*Job, error)
	GetExecutions(jobName string) ([]*Execution, error)
	GetLastExecutionGroup(jobName string) ([]*Execution, error)
	GetExecutionGroup(execution *Execution) ([]*Execution, error)
	GetGroupedExecutions(jobName string) (map[int64][]*Execution, []int64, error)
	Shutdown() error
	Snapshot(w io.WriteCloser) error
	Restore(r io.ReadCloser) error
}
