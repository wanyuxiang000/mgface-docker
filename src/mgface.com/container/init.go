package container

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
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
		return fmt.Errorf("获取用户输入的指令为空,不能为空!")
	}

	setUpContainerMount()

	path, err := exec.LookPath(cmdArray[0])

	if err != nil {
		logrus.Errorf("没有定位到执行的命令: %v", err)
		return err
	}
	logrus.Infof("找到命令的执行路径: %s", path)
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
	go func() {
		http.ListenAndServe(":4444", nil)
	}()
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

func setUpContainerMount() {
	pwd, _ := os.Getwd()
	logrus.Infof("当前的location: %s", pwd)
	//makes the mount namespace works properly on my archlinux computer, systemd made "/" mounted as shared by default.
	//systemd加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示声明你要这个新的mount namespace独立
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	err := pivotRoot(pwd)
	logrus.Infof("pivotRoot切换->%v", err)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		fmt.Println("mount proc 发生错误:", err)
	}
	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755"); err != nil {
		fmt.Println("mount tmpfs 发生错误:", err)
	}
}

func pivotRoot(root string) error {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_PRIVATE|syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("挂载rootfs给自己发生错误:%v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	fmt.Println("pivotDir->", pivotDir)
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	//pivotRoot把当前进程的root系统移动到putold文件夹，然后让new_root成为新root的文件系统
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	//把当前root表示的目录切换为根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("umount pivot_root dir %v", err)
	}
	return os.Remove(pivotDir)
}
