package container

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"mgface.com/containerInfo"
	"strconv"
	"syscall"
)

func StopContainer(containerName string) error {
	containerURL := fmt.Sprintf(containerInfo.DefaultInfoLocation, containerName)
	configURL := containerURL + containerInfo.ConfigName
	logrus.Infof("开始停止容器%s.配置文件位置:%s", containerName, configURL)
	content, _ := ioutil.ReadFile(configURL)
	var containerinfo containerInfo.ContainerInfo
	json.Unmarshal(content, &containerinfo)
	pid := containerinfo.Pid
	ipid, _ := strconv.Atoi(pid)
	if err := syscall.Kill(ipid, syscall.SIGTERM); err != nil {
		logrus.Errorf("中断进程%d失败,异常信息为:%v", ipid, err)
	}
	containerinfo.Status = containerInfo.STOP
	containerinfo.Pid = ""
	content, _ = json.Marshal(containerinfo)
	if err := ioutil.WriteFile(configURL, content, 0622); err != nil {
		logrus.Errorf("写文件%s失败，错误日志:%v", configURL, err)
	}
	return nil
}
