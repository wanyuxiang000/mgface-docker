package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"mgface.com/containerInfo"
	"strconv"
	"syscall"
	"time"
)

func StartContainer(containerName string) error {
	containerURL := fmt.Sprintf(containerInfo.DefaultInfoLocation, containerName)
	configURL := containerURL + containerInfo.ConfigName
	logrus.Infof("开始启动一个存在的容器%s.配置文件位置:%s", containerName, configURL)
	content, _ := ioutil.ReadFile(configURL)
	var containerinfo containerInfo.ContainerInfo
	json.Unmarshal(content, &containerinfo)

	if containerinfo.Status != containerInfo.STOP {
		return errors.New("启动容器状态应该为stopped.")
	}

	pid := containerinfo.Pid
	ipid, _ := strconv.Atoi(pid)

	if err := syscall.Kill(ipid, syscall.SIGCONT); err != nil {
		return errors.New(fmt.Sprintf("重新开始一个停止的进程%d失败,异常信息为:%v", ipid, err))
	}
	containerinfo.Status = containerInfo.RUNNING
	containerinfo.StoppedTime = time.Now().Format("2006-01-02 15:04:05")
	content, _ = json.MarshalIndent(containerinfo, "", "   ") //美化输出缩进格式
	if err := ioutil.WriteFile(configURL, content, 0622); err != nil {
		logrus.Errorf("写文件%s失败，错误日志:%v", configURL, err)
	}
	return nil
}
