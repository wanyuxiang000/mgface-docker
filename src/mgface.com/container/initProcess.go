package container

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	logrus.Infof("接收到的命令:%s", cmdArray)
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error,cmdArray is nil")
	}

	setUpMount()

	path, err := exec.LookPath(cmdArray[0])

	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}
	logrus.Infof("Find path %s", path)
	//MS_NOEXEC 本文件系统不允许运行其他程序
	//MS_NOSUID 不允许setuserId和setGroupId
	//MS NODEV 这个参数是自 从 Linux 2.4 以来，所有 mount 的系统都会默认设定的参数
	//defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	//syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	//argv := []string{command}
	//这个方法是调用内核的int execve(const char *filename,char * const argv[],char *const envp[])
	//它的作用是执行当前filename对应的程序。它会覆盖当前进程的镜像、数据和堆械等信息，包括 PID.这些都会被将要运行的进程覆盖掉。
	//也就是说，调用这个方法，将用户指定的进程运行起来，把最初的 init 进程给替换掉，这样当进入到容器内部的时候，就会发现容器内的第一个程序就是我们指定的进程了
	logrus.Infof("命令行:%s", cmdArray[0:])
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, _ := ioutil.ReadAll(pipe)
	msgstr := string(msg)
	return strings.Split(msgstr, " ")
}

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_PRIVATE|syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("挂载rootfs给自己发生错误:%v", err)
	}
	pivotDir := filepath.Join(root, ".pivot_root")
	fmt.Println("pivotDir->", pivotDir)
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	//pivotRoot把当前进程的root系统移动到putold文件夹，然后让new_root成为新root的文件系统
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")

	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("umount pivot_root dir %v", err)
	}

	return os.Remove(pivotDir)

}

func setUpMount() {
	pwd, _ := os.Getwd()
	logrus.Infof("当前的location: %s", pwd)
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	err := pivotRoot(pwd)
	logrus.Infof("pivotRoot切换->%v", err)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		fmt.Println("mount proc 发生错误:", err)
	}
	if err := syscall.Mount("tmpfs", "/dev", "tempfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755"); err != nil {
		fmt.Println("mount tmpfs 发生错误:", err)
	}
}
