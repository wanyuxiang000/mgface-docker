package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cgroupMemoryMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		fmt.Printf("当前Pid %d \n", syscall.Getpid())
		cmd := exec.Command("sh", "-c", "stress --vm-bytes 200m --vm-keep -m 1")
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		log.Fatal("发生错误:%s", err)
		os.Exit(-1)
	} else {
		log.Printf("fork出来进程的PID(外部命名空间)：%d", cmd.Process.Pid)
		err:=os.Mkdir(path.Join(cgroupMemoryMount, "testmemorylimit"), 0755)
		if err!=nil {
			log.Fatal("错误的创建文件,发生错误.%s",err)
			os.Exit(-1)
		}
		//将容器进程加入到这个cgroup
		ioutil.WriteFile(path.Join(cgroupMemoryMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		ioutil.WriteFile(path.Join(cgroupMemoryMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
	}
	cmd.Process.Wait()
}
