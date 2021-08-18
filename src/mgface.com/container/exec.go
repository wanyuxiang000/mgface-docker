package container

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mgface.com/constVar"
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
	dirURL := fmt.Sprintf(constVar.DefaultInfoLocation, containerName)
	config := dirURL + constVar.ConfigName
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
	//获取对应的PID环境变量,其实也是容器的环境变量
	containerEnvs := getEnvsByPid(pid)
	logrus.Infof("获取容器PID:%d的环境变量信息:%v", pid, containerEnvs)
	//将宿主机的环境变量和容器的环境变量都放置到exec进程内
	cmd.Env = append(os.Environ(), containerEnvs...)

	if err := cmd.Run(); err != nil {
		logrus.Errorf("执行容器:%s，发生异常:%v", containerName, err)
	}
}

func getEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Errorf("读取文件 %s 发生异常: %v", path, err)
		return nil
	}
	//多个环境变量中的分隔符是\u0000
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}
