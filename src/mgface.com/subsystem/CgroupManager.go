package subsystem

type CgroupManager struct {
	Path     string
	Resource *ResouceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int) {
	for _, sub := range SybsystemsIns {
		sub.Apply(c.Path, pid)
	}
}

func (c *CgroupManager) Set(res *ResouceConfig){
	for _, sub := range SybsystemsIns {
		sub.Set(c.Path, res)
	}
}

func (c *CgroupManager) Destory(){
	for _, sub := range SybsystemsIns {
		sub.Remove(c.Path)
	}
}
