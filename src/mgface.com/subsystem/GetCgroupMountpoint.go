package subsystem

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

func FindCgroupMountpoint(subsystem string) string {
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

func GetCgroupPath(subsystem string, cgrouppath string, autocreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgrouppath)); err == nil || (autocreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgrouppath), 0755); err != nil {
				return "", fmt.Errorf("错误的创建了cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgrouppath), nil
	} else {
		return "", fmt.Errorf("错误的创建了cgroup %v", err)
	}
}
