package main

import (
	"github.com/victorcoder/dkron/dkron"
)

// FilesOutput plugin that saves each execution log
// in it's own file in the file system.
type Shell struct {
	Param1 string
	Param2 bool
}

// Process method of the plugin
func (s *Shell) Execute(args *dkron.ExecuteRequest) ([]byte, error) {
	return []byte("Response from plugin"), nil
}
