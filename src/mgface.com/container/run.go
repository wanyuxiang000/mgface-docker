package container

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"mgface.com/aufs"
	"mgface.com/cgroup"
	"mgface.com/constVar"
	"mgface.com/containerInfo"
	"mgface.com/containerNet"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func newParentProcess(tty bool, volume string, containerName string, envs []string) (*exec.Cmd, *os.File) {
	r, w, _ := os.Pipe()
	args := []string{"init"}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC, // | syscall.CLONE_NEWUSER,
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
	} else {
		//创建输出日志的目录和文件
		dirURL := fmt.Sprintf(constVar.DefaultInfoLocation, containerName)
		os.MkdirAll(dirURL, 0622)
		stdLogFile := dirURL + constVar.ContainerLog
		stdout, _ := os.Create(stdLogFile)
		logrus.Infof("生成容器:%s的日志文件:%s", containerName, stdLogFile)
		cmd.Stdout = stdout
	}
	// 外带着这个文件句柄去创建子进程因为1个进程默认会有3个文件描述符,分别是标准输入、标准输出、标准错误.
	// 这3个是子进程一创建就会默认带着的,那么外带的这个文件描述符就成为了第4个
	cmd.ExtraFiles = []*os.File{r}

	//设置cmd的目录
	cmd.Dir = fmt.Sprintf(constVar.Cmd, containerName)

	//设置环境变量
	//默认情况下,新启动进程的环境变量都是继承于原来父进程的,但是如果手动指定了环境变量,那么就会覆盖掉原来继承自父进程的变量。由于在容器的进程中，
	//有时候还需要使用原来父进程的环境变量,比如PATH等,因此这里会使用os.Environ()来获取宿主机的环境变量,然后把自定义的变量加进去。
	cmd.Env = append(os.Environ(), envs...)

	//设置好容器进程的挂载点(作为容器的文件系统)
	aufs.NewFileSystem(volume, containerName)
	return cmd, w
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("所有的命令:%s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func RunContainer(tty bool, command []string, res *cgroup.ResouceConfig, volume string, containerName string, envs []string, network string, portMapping []string) {
	//获取容器名称
	containerName, id := containerInfo.GetContainerName(containerName)
	//当前进程创建容器的父进程
	parent, writePipe := newParentProcess(tty, volume, containerName, envs)
	if err := parent.Start(); err != nil {
		logrus.Fatal("发生错误:%s", err)
	}
	//获取容器的PID
	pid := parent.Process.Pid
	//记录容器信息
	containerInfo.RecordContainerInfo(pid, command, containerName, id, volume)

	//设置Cgroup
	cgroup.SetCgroup(fmt.Sprintf(constVar.CgroupName, containerName), res, pid)

	//设置容器连接网络
	if network != "" {
		//初始化容器网络和驱动
		containerNet.InitNetworkAndNetdriver()

		containerInfo := &containerInfo.ContainerInfo{
			Id:          id,
			Pid:         strconv.Itoa(pid),
			Name:        containerName,
			PortMapping: portMapping,
		}
		logrus.Info("**************开始配置网络**************")
		if err:=containerNet.Connect(network, containerInfo);err!=nil{
			logrus.Fatal("配置网络发生错误:%s", err.Error())
			os.Exit(-1)
		}
		logrus.Info("**************结束[end]配置网络**************")
	}

	//向容器进程进行通信
	sendInitCommand(command, writePipe)

	//假如启用了tty的话，那么父类进程等待子类进程结束
	if tty {
		parent.Wait()
		logrus.Infof("退出当前进程:%d,时间为:%s", os.Getpid(), time.Now().Format("2006-01-02 15:04:05"))
		logrus.Infof("开始清理环境...")
		//删除容器信息
		containerInfo.DeleteContainerInfo(containerName)
		//删除挂载点数据
		aufs.DeleteFileSystem(volume, containerName)
	} else {
		logrus.Info("等待3秒为了初始化宿主机监听端口信息等...")
		time.Sleep(3*time.Second)
		logrus.Infof("不启用tty,父进程直接运行完毕,子进程进行detach分离给操作系统的init托管.")
	}
}
