package impl

import (
	"fmt"
	"io/ioutil"
	"mgface.com/cgroup"
	"os"
	"path"
	"strconv"
)

type CpuSubSyetem struct {
}

func (s *CpuSubSyetem) Name() string {
	return "cpu"
}
func (s *CpuSubSyetem) Set(cgroupPath string, res *cgroup.ResouceConfig) error {
	if subsysCgroupPath, err := cgroup.GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuShare != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
				return fmt.Errorf("设置使用cpu权重失败 %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}
func (s *CpuSubSyetem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := cgroup.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("往%s的task添加进程ID失败. %v", subsysCgroupPath, err)
		}
		return nil
	} else {
		return fmt.Errorf("往%s的task添加进程ID失败. %v", subsysCgroupPath, err)
	}
}
func (s *CpuSubSyetem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := cgroup.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}
