package dcron

import (
	"os"
	"os/exec"
)

// spawn command that specified as proc.
func spawnProc(proc string) (*exec.Cmd, error) {
	cs := []string{"/bin/bash", "-c", proc}
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ())

	log.Printf("Starting %s\n", proc)
	err := cmd.Start()
	if err != nil {
		log.Errorf("Failed to start %s: %s\n", proc, err)
		return nil, err
	}
	return cmd, nil
}
