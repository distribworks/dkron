// +build !windows

package main

import (
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func setCmdAttr(cmd *exec.Cmd, config map[string]string) error {
	su, _ := config["su"]
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
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		}
	}
	return nil
}
