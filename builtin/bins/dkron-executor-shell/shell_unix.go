// +build !windows

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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

func calculateMemory(pid int) (uint64, error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/smaps", pid))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	res := uint64(0)
	rfx := []byte("Rss:")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, rfx) {
			var size uint64
			_, err := fmt.Sscanf(string(line[4:]), "%d", &size)
			if err != nil {
				return 0, err
			}
			res += size
		}
	}
	if err := r.Err(); err != nil {
		return 0, err
	}
	return res, nil
}
