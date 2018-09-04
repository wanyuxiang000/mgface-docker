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
	//这个方法是调用内核的int execve(const char *filename,char * const argv[],char *const envp[])
	//它的作用是执行当前filename对应的程序。它会覆盖当前进程的镜像、数据和堆械等信息，包括 PID.这些都会被将要运行的进程覆盖掉。
	//也就是说，调用这个方法，将用户指定的进程运行起来，把最初的 init 进程给替换掉，这样当进入到容器内部的时候，就会发现容器内的第一个程序就是我们指定的进程了
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}