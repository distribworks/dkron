package dkron

type ExecutionProcessor interface {
	Process(execution *Execution) *Execution
}
