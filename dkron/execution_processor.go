package dkron

type ExecutionProcessor interface {
	Process(args *ExecutionProcessorArgs) Execution
}

type ExecutionProcessorArgs struct {
	Execution Execution
	Config    map[string]string
}
