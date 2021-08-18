package cgroup

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

//设置cgroup参数
func SetCgroup(cgroupName string, res *ResouceConfig, pid int) *cgroupManager {
	manager := newCgroupManager(cgroupName)
	defer func() {
		if err := recover(); err != nil {
			logrus.Info("发生了panic异常,信息为:%s,进行资源销毁.", err)
			manager.Destory()
		}
	}()
	//设置资源限制
	manager.Set(res)
	logrus.Infof("设置cgroup资源限制完成....")
	//将容器进程加入到各个cgroup subsystem
	manager.Apply(pid)
	logrus.Infof("应用cgroup添加PID进相应限制文件的Task文件完成...")
	return manager
}

func findCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

func getCgroupPath(subsystem string, cgroupPath string, autocreate bool) (string, error) {
	cgroupRoot := findCgroupMountpoint(subsystem)
	cgroupURL := path.Join(cgroupRoot, cgroupPath)
	logrus.Infof("子系统%s 创建cgroup路径:%s", subsystem, cgroupURL)
	if _, err := os.Stat(cgroupURL); err == nil || (autocreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(cgroupURL, 0755); err != nil {
				return "", fmt.Errorf("错误的创建了cgroup %v", err)
			}
		}
		return cgroupURL, nil
	} else {
		return "", fmt.Errorf("错误的创建了cgroup %v", err)
	}
}
