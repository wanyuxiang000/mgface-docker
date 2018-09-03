package main

import (
	logs "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	user,err:=user.Lookup("nobody")
	if err == nil {
		logs.Errorf("uid=%s,gid=%s",user.Uid,user.Gid)
	}
	uid,_:=strconv.Atoi(user.Uid)
	gid,_:=strconv.Atoi(user.Gid)

	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logs.Errorf("发生了致命错误:%s", err)

	}

	os.Exit(-1)
}
