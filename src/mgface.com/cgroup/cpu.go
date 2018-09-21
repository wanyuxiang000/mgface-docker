package cgroup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type cpuSubSyetem struct {
}

func (s *cpuSubSyetem) Name() string {
	return "cpu"
}
func (s *cpuSubSyetem) Set(cgroupPath string, res *ResouceConfig) error {
	if subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuShare != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
				//panic("设置使用cpu权重失败 %v")
				return fmt.Errorf("设置使用cpu权重失败 %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}
func (s *cpuSubSyetem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("往%s的task添加进程ID失败. %v", subsysCgroupPath, err)
		}
		return nil
	} else {
		return fmt.Errorf("往%s的task添加进程ID失败. %v", subsysCgroupPath, err)
	}
}
func (s *cpuSubSyetem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}
