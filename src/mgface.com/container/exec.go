package container

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"mgface.com/containerInfo"
	_ "mgface.com/nsenter"
	"os"
	"os/exec"
	"strings"
)

const (
	EnvExecPid = "mgfacedocker_pid"
	EnvExecCmd = "mgfacedocker_cmd"
)

func getContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(containerInfo.DefaultInfoLocation, containerName)
	config := dirURL + containerInfo.ConfigName
	content, _ := ioutil.ReadFile(config)
	var info containerInfo.ContainerInfo
	json.Unmarshal(content, &info)
	return info.Pid, nil
}

func ExecContainer(containerName string, cmdArray []string) {
	pid, _ := getContainerPidByName(containerName)
	cmdStr := strings.Join(cmdArray, " ")
	logrus.Infof("容器ID:%s", pid)
	logrus.Infof("command:%s", cmdStr)
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	os.Setenv(EnvExecCmd, cmdStr)
	os.Setenv(EnvExecPid, pid)

	if err := cmd.Run(); err != nil {
		logrus.Errorf("执行容器:%s，发生异常:%v", containerName, err)
	}
}
