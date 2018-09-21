package cgroup

import (
	"mgface.com/cgroup.limit"
)

type ResouceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type cgroupSubsytem interface {
	Name() string
	Set(path string, res *ResouceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

var cgroupSubsytems = []cgroupSubsytem{
	&cgroup_limit.MemorySubSyetem{},
	&cgroup_limit.CpuSubSyetem{},
	&cgroup_limit.CpusetSubSystem{},
}
