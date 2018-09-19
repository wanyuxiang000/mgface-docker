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
	ENV_EXEC_PID = "mydocker_id"
	ENV_EXEC_CMD = "mydocker_cmd"
)

func GetContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(containerInfo.DefaultInfoLocation, containerName)
	config := dirURL + containerInfo.ConfigName
	content, _ := ioutil.ReadFile(config)
	var info containerInfo.ContainerInfo
	json.Unmarshal(content, &info)
	return info.Pid, nil
}

func ExecContainer(containerName string, cmdArray []string) {
	pid, _ := GetContainerPidByName(containerName)
	cmdStr := strings.Join(cmdArray, " ")
	logrus.Infof("容器ID:%s", pid)
	logrus.Infof("command:%s", cmdStr)
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	os.Setenv(ENV_EXEC_CMD, cmdStr)
	os.Setenv(ENV_EXEC_PID, pid)

	if err := cmd.Run(); err != nil {
		logrus.Errorf("执行容器:%s，发生异常:%v", containerName, err)
	}
}
