//go:build !windows
// +build !windows

package shell

import (
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func setCmdAttr(cmd *exec.Cmd, config map[string]string) error {
	su := config["su"]
	if su != "" {
		var uid, gid int
		parts := strings.Split(su, ":")
		u, err := user.Lookup(parts[0])
		if err != nil {
			return err
		}
		uid, _ = strconv.Atoi(u.Uid)
		if len(parts) > 1 {
			g, err := user.LookupGroup(parts[1])
			if err != nil {
				return err
			}
			gid, _ = strconv.Atoi(g.Gid)
		} else {
			gid, _ = strconv.Atoi(u.Gid)
		}
		cmd.SysProcAttr.Credential = &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		}
	}

	jobTimeout := config["timeout"]
	if jobTimeout != "" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
	return nil
}

func processKill(cmd *exec.Cmd) error {
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) // note the minus sign
}
