package container

import (
	"github.com/Sirupsen/logrus"
	"mgface.com/subsystem"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER | syscall.CLONE_NEWIPC,
	}

	//相当于 cmd 进程认为自己是以 root 执行的，但其实最终的操作受制于 宿主机0(root)这个用户

	//设置该进程在新命名空间中以root用户执行。而这个root用户则是映射到host上用户id为 0(root)、组 id 为 0(root) 的用户
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: 0,
		Gid: 0,
	}

	cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: 0, Size: 1}}
	cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: 0, Size: 1}}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func Run(tty bool, command string,res *subsystem.ResouceConfig) {
	cmd := NewParentProcess(tty, command)
	if err := cmd.Start(); err != nil {
		logrus.Fatal("发生错误:%s", err)
	}
	manager:=subsystem.NewCgroupManager("mgface-cgroup")
	defer manager.Destory()
	//设置资源限制
	manager.Set(res)
	//将容器进程加入到各个subsystem挂载对于的cgroup
	manager.Apply(cmd.Process.Pid)
	cmd.Wait()
	os.Exit(-1)
}
