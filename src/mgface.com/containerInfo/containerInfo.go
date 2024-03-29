package containerInfo

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mgface.com/constVar"
	"os"
	"strconv"
	"strings"
	"time"
)

type ContainerInfo struct {
	Pid         string   `json:"pid"`         //容器init进程在宿主机的PID
	Id          string   `json:"id"`          //容器ID
	Name        string   `json:"name"`        //容器名称
	Command     string   `json:"command"`     //容器内init进程的执行命令
	CreatedTime string   `json:"createdTime"` //创建时间
	StoppedTime string   `json:"stoppedTime"` //停止时间
	Status      string   `json:"status"`      //容器状态
	Volume      string   `json:"volume"`      //挂载卷
	PortMapping []string `json:"portmapping"` //端口映射
}

func GetContainerName(containerName string) (string, string) {
	id := randStringBuffer(10)
	if containerName == "" {
		containerName = id
	}
	return containerName, id
}

func DeleteContainerInfo(containerName string) {
	dirURL := fmt.Sprintf(constVar.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		logrus.Errorf("删除目录 %s 失败:%v", dirURL, err)
	}
}

func GetContainerInfo(file os.FileInfo) (*ContainerInfo, error) {
	containerName := file.Name()
	configFileDir := fmt.Sprintf(constVar.DefaultInfoLocation, containerName)
	configFileDir = configFileDir + constVar.ConfigName
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		logrus.Errorf("读取文件失败%v", err)
	}
	var containerInfo ContainerInfo
	json.Unmarshal(content, &containerInfo)
	return &containerInfo, nil
}

//记录容器信息
func RecordContainerInfo(containerPID int, commandArray []string, containerName string, id string, volume string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, ",")

	containerInfo := &ContainerInfo{
		Pid:         strconv.Itoa(containerPID),
		Id:          id,
		Name:        containerName,
		Command:     command,
		CreatedTime: createTime,
		StoppedTime: "",
		Status:      constVar.RUNNING,
		Volume:      volume,
	}

	dirUrl := fmt.Sprintf(constVar.DefaultInfoLocation, containerName)
	os.MkdirAll(dirUrl, 0622)
	fileName := dirUrl + "/" + constVar.ConfigName
	file, _ := os.Create(fileName)
	defer file.Close()
	jsonBytes, _ := json.MarshalIndent(containerInfo, "", "   ") //美化输出缩进格式
	content := append(jsonBytes, []byte("\n")...)
	file.Write(content)
	return containerName, nil
}
