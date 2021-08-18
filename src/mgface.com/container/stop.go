package container

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mgface.com/constVar"
	"mgface.com/containerInfo"
	"strconv"
	"syscall"
	"time"
)

func StopContainer(containerName string) error {
	containerURL := fmt.Sprintf(constVar.DefaultInfoLocation, containerName)
	configURL := containerURL + constVar.ConfigName
	logrus.Infof("开始停止容器%s.配置文件位置:%s", containerName, configURL)
	content, _ := ioutil.ReadFile(configURL)
	var containerinfo containerInfo.ContainerInfo
	json.Unmarshal(content, &containerinfo)
	pid := containerinfo.Pid
	ipid, _ := strconv.Atoi(pid)
	if err := syscall.Kill(ipid, syscall.SIGSTOP); err != nil {
		logrus.Errorf("停止进程%d失败,异常信息为:%v", ipid, err)
	}
	containerinfo.Status = constVar.STOP
	containerinfo.StoppedTime = time.Now().Format("2006-01-02 15:04:05")
	//containerinfo.Pid = ""
	content, _ = json.MarshalIndent(containerinfo, "", "   ") //美化输出缩进格式
	content = append(content, []byte("\n")...)
	if err := ioutil.WriteFile(configURL, content, 0622); err != nil {
		logrus.Errorf("写文件%s失败，错误日志:%v", configURL, err)
	}
	return nil
}
