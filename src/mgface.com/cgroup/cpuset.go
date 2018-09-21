package cgroup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type cpusetSubSystem struct {
}

func (s *cpusetSubSystem) Set(cgroupPath string, res *ResouceConfig) error {
	if subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuSet != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("设置绑定进程到指定CPU Core失败,%v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *cpusetSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

func (s *cpusetSubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("往 %s的task添加进程ID失败. %v", subsysCgroupPath, err)
		}
		return nil
	} else {
		return fmt.Errorf("往 %s的task添加进程ID失败. %v", subsysCgroupPath, err)
	}
}

func (s *cpusetSubSystem) Name() string {
	return "cpuset"
}
