package subsystem

type ResouceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsytem interface {
	Name() string
	Set(path string, res *ResouceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

var (
	SybsystemsIns = []Subsytem{
		&MemorySubSyetem{},
		&CpuSubSyetem{},
		&CpusetSubSystem{},
	}
)
