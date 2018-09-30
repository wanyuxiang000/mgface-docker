package cgroup

type cgroupManager struct {
	Path     string
	Resource *ResouceConfig
}

func newCgroupManager(path string) *cgroupManager {
	return &cgroupManager{
		Path: path,
	}
}

func (c *cgroupManager) Apply(pid int) {
	for _, sub := range cgroupSubsytems {
		sub.Apply(c.Path, pid)
	}
}

func (c *cgroupManager) Set(res *ResouceConfig) {
	for _, sub := range cgroupSubsytems {
		sub.Set(c.Path, res)
	}
}

func (c *cgroupManager) Destory() {
	for _, sub := range cgroupSubsytems {
		sub.Remove(c.Path)
	}
}
