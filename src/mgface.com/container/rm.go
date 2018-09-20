package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mgface.com/containerInfo"
	"os"
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

	//TODO 删除挂载的文件
	//删除当前目录
	return os.RemoveAll(dirURL)
}