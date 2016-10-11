package dkron

type Outputter interface {
	Output(execution *Execution) string
}
