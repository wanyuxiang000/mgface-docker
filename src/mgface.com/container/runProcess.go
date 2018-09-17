package container

import (
	"github.com/Sirupsen/logrus"
	"mgface.com/aufs"
	"mgface.com/subsystem"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func NewParentProcess(tty bool,volume string) (*exec.Cmd, *os.File) {
	r, w, _ := os.Pipe()
	args := []string{"init"}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,// | syscall.CLONE_NEWUSER,
	}


	//相当于 cmd 进程认为自己是以 root 执行的，但其实最终的操作受制于 宿主机0(root)这个用户

	//设置该进程在新命名空间中以root用户执行。而这个root用户则是映射到host上用户id为 0(root)、组 id 为 0(root) 的用户
	//cmd.SysProcAttr.Credential = &syscall.Credential{
	//	Uid: 0,
	//	Gid: 0,
	//}
	//
	//cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: 0, Size: 1}}
	//cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{{ContainerID: 0, HostID: 0, Size: 1}}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	//会外带着这个文件句柄去创建子进程
	//因为 1 个进程默认会有 3 个文件描述符，分别是标准输入、标准输出、标准错误。这3个
	//是子进程一创建的时候就会默认带着的，那么外带的这个文件描述符理所当然地就成为了第4个
	cmd.ExtraFiles = []*os.File{r}

	//设置cmd的目录
	cmd.Dir = "/root/mnt"
	aufs.NewWorkspace("/root","/root/mnt",volume)
	return cmd, w
}

func Run(tty bool, command []string, res *subsystem.ResouceConfig,volume string) {
	parent, writePipe := NewParentProcess(tty,volume)
	if err := parent.Start(); err != nil {
		logrus.Fatal("发生错误:%s", err)
	}
	manager := subsystem.NewCgroupManager("mgface-cgroup")
	defer manager.Destory()
	//设置资源限制
	manager.Set(res)
	//将容器进程加入到各个subsystem挂载对于的cgroup
	manager.Apply(parent.Process.Pid)

	sendInitCommand(command, writePipe)
	//假如启用了tty的话，那么父类进程等待子类进程结束
	if tty {
		parent.Wait()
	}
	//logrus.Infof("退出当前进程:%s",time.Now().Format("2006-01-02 15:04:05"))
	//logrus.Infof("开始清理环境...")
	//aufs.DeleteWorkSpace("/root","/root/mnt",volume)
	//os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("所有的命令:%s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
