package aufs

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"mgface.com/constVar"
	"os"
	"os/exec"
	"strings"
)

func NewFileSystem(volume string, containerName string, command []string) {
	fmt.Println("解析的command:",command)
	logrus.Infof("1)创建只读层...")
	createReadOnlyLayer(containerName,command)
	logrus.Infof("2)创建可写层...")
	createWriteLayer(containerName)
	logrus.Infof("3)创建挂载点...")
	createMountPoint(containerName)
	logrus.Infof("4)挂载卷映射...")
	volumeMapping(volume, containerName)
}

//创建只读层
func createReadOnlyLayer(containerName string,command []string) {
	busyboxUrl := fmt.Sprintf(constVar.FileSystemURL,command[0], containerName)
	busyboxTarURL := fmt.Sprintf(constVar.FileSystemTarURL,command[0])
	exist, _ := pathExit(busyboxUrl)
	if exist == false {
		if err := os.MkdirAll(busyboxUrl, 0777); err != nil {
			logrus.Errorf("创建目录%s  发生异常%v", busyboxUrl, err)
		}
		logrus.Infof("准备解压tar包:%s  到目录:%s", busyboxTarURL, busyboxUrl)
		exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxUrl).CombinedOutput()
	} else {
		logrus.Info("挂载的文件系统存在.目录:%s", busyboxUrl)
		fileInfos, _ := ioutil.ReadDir(busyboxUrl)
		if len(fileInfos) < 2 {
			os.RemoveAll(busyboxUrl)
			logrus.Errorf("文件系统目录下面不存在文件,解压tar文件系统 %s 到 %s .", busyboxTarURL, busyboxUrl)
			createReadOnlyLayer(containerName,command)
		}
	}
}

func createWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(constVar.WriteLayer, containerName)
	os.MkdirAll(writeURL, 0777)
}

func createMountPoint(containerName string) error {
	containerNameURL := fmt.Sprintf(constVar.MntURL, containerName)
	logrus.Infof("开始创建挂载点目录%s.", containerNameURL)
	os.MkdirAll(containerNameURL, 0777)
	cmd := exec.Command("mount", "-t", "aufs", "-o", fmt.Sprintf(constVar.MountAufsDirs, containerName, containerName), "none", containerNameURL)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	logrus.Info("开始执行系统挂载.")
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("创建挂载点发生严重错误:%v", err)
		return errors.New(err.Error())
	}
	return nil
}

func pathExit(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}

//挂载卷映射
func volumeMapping(volumes string, containerName string) {
	if volumes != "" {
		volumeUrls := strings.Split(volumes, ",")
		for _, v := range volumeUrls {
			logrus.Infof("挂载  [%s]  挂载点", v)
			volume := strings.Split(v, ":")
			if len(volume) == 2 && volume[0] != "" && volume[1] != "" {
				mountVolume(volume, containerName)
			} else {
				logrus.Errorf("volume参数配置错误,请使用v1:v2这样的格式.")
			}
		}
	} else {
		logrus.Info("没有挂载卷数据.")
	}
}

func mountVolume(volumeUrls []string, containerName string) {
	parentUrl := volumeUrls[0]
	os.MkdirAll(parentUrl, 0777)

	containerUrl := volumeUrls[1]
	containerVolumeURL := fmt.Sprintf(constVar.MntURL, containerName) + containerUrl
	os.MkdirAll(containerVolumeURL, 0777)

	//把宿主机文件挂到容器挂载点上
	dirs := "dirs=" + parentUrl
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount volume 错误.%v", err)
	}
}
