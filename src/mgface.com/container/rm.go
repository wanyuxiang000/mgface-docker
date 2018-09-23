package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mgface.com/aufs"
	"mgface.com/containerInfo"
	"strconv"
	"syscall"
)

func RemoveContainer(containerName string) error {
	dirURL := fmt.Sprintf(containerInfo.DefaultInfoLocation, containerName)
	configURL := dirURL + containerInfo.ConfigName
	content, _ := ioutil.ReadFile(configURL)
	var cinfo containerInfo.ContainerInfo
	json.Unmarshal(content, &cinfo)
	if cinfo.Status != containerInfo.STOP {
		//logrus.Errorf("不能删除容器状态不为stopped,请先执行stop指令再删除.")
		return errors.New("不能删除容器状态不为stopped,请先执行stop指令再删除.")
	}
	//暂停当前进程信息
	pid := cinfo.Pid
	ipid, _ := strconv.Atoi(pid)

	if err := syscall.Kill(ipid, syscall.SIGTERM); err != nil {
		return errors.New(fmt.Sprintf("杀掉进程%d失败,异常信息为:%v", ipid, err))
	}

	//删除当前目录
	containerInfo.DeleteContainerInfo(containerName)
	//删除挂载点数据
	aufs.DeleteFileSystem(cinfo.Volume, containerName)
	return nil
}
