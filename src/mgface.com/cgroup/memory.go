package cgroup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type memorySubSyetem struct {
}

func (s *memorySubSyetem) Name() string {
	return "memory"

}
func (s *memorySubSyetem) Set(cgroupPath string, res *ResouceConfig) error {
	if sub, err := getCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(sub, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("设置memory失败,%v", err)
			}
		}
	}
	return nil
}
func (s *memorySubSyetem) Apply(cgroupPath string, pid int) error {
	if sub, err := getCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(sub, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("把PID：%d 添加到%s task文件失败", pid, sub)
		}
	} else {
		return fmt.Errorf("把PID：%d 添加到%s task文件失败", pid, sub)
	}
	return nil

}
func (s *memorySubSyetem) Remove(cgroupPath string) error {
	if sub, err := getCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.Remove(sub)
	} else {
		return err
	}
}
