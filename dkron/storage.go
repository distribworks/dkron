package dkron

type Storage interface {
	SetJob(job *Job, copyDependentJobs bool) error
	DeleteJob(name string) (*Job, error)
	SetExecution(execution *Execution) (string, error)
	DeleteExecutions(jobName string) error
	SetExecutionDone(execution *Execution) (bool, error)

	GetJobs(options *JobOptions) ([]*Job, error)
	GetJob(name string, options *JobOptions) (*Job, error)
	GetExecutions(jobName string) ([]*Execution, error)
	GetLastExecutionGroup(jobName string) ([]*Execution, error)
	GetExecutionGroup(execution *Execution) ([]*Execution, error)
	GetGroupedExecutions(jobName string) (map[int64][]*Execution, []int64, error)
	Shutdown() error
}
