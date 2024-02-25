//go:build windows
// +build windows

package shell

import (
	"os/exec"
)

func setCmdAttr(cmd *exec.Cmd, config map[string]string) error {
	return nil
}

func processKill(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}
