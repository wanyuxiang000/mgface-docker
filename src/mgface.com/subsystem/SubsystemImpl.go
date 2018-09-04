package subsystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSyetem struct {
}

func (s *MemorySubSyetem) Name() string {
	return "memory"

}
func (s *MemorySubSyetem) Set(cgroupPath string, res *ResouceConfig) error {
	if sub, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(sub, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("设置cgroup内存失败,%v", err)
			}
		}
	}
	return nil
}
func (s *MemorySubSyetem) Apply(cgroupPath string, pid int) error {
	if sub, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(sub, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("把PID：%d 添加到task文件失败", pid)
		}
	} else {
		return fmt.Errorf("把PID：%d 添加到task文件失败", pid)
	}
	return nil

}
func (s *MemorySubSyetem) Remove(cgroupPath string) error {
	if sub, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.Remove(sub)
	} else {
		return err
	}
}
