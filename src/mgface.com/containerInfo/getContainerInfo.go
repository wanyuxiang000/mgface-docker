package containerInfo

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
)

func GetContainerInfo(file os.FileInfo) (*ContainerInfo, error) {
	containerName := file.Name()
	configFileDir := fmt.Sprintf(DefaultInfoLocation, containerName)
	configFileDir = configFileDir + ConfigName
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		logrus.Errorf("读取文件失败%v", err)
	}
	var containerInfo ContainerInfo
	json.Unmarshal(content, &containerInfo)
	return &containerInfo, nil
}
