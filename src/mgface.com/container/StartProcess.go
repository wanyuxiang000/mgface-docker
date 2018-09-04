package container

import (
	"github.com/Sirupsen/logrus"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, _ []string) error {
	logrus.Infof("接收到的命令:%s", command)
	//MS_NOEXEC 本文件系统不允许运行其他程序
	//MS_NOSUID 不允许setuserId和setGroupId
	//MS NODEV 这个参数是自 从 Linux 2.4 以来，所有 mount 的系统都会默认设定的参数
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}
