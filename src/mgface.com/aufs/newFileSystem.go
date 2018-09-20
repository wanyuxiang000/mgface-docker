package aufs

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"mgface.com/constVar"
	"os"
	"os/exec"
	"strings"
)

func NewFileSystem(volume string) {
	logrus.Infof("1)创建只读层...")
	createReadOnlyLayer()
	logrus.Infof("2)创建可写层...")
	createWriteLayer()
	logrus.Infof("3)创建挂载点...")
	createMountPoint()
	logrus.Infof("4)挂载卷映射...")
	volumeMapping(volume)
}

//创建只读层
func createReadOnlyLayer() {
	busyboxUrl := constVar.FileSystemURL
	busyboxTarURL := constVar.FileSystemTarURL
	exist, _ := pathExit(busyboxUrl)
	if exist == false {
		if err := os.Mkdir(busyboxUrl, 0777); err != nil {
			logrus.Errorf("创建目录%s  发生异常%v", busyboxUrl, err)
		}
		logrus.Infof("准备解压tar包:%s  到目录:%s", busyboxTarURL, busyboxUrl)
		exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxUrl).CombinedOutput()
	} else {
		fileInfos, _ := ioutil.ReadDir(busyboxUrl)
		if len(fileInfos) < 1 {
			logrus.Errorf("文件系统目录下面不存在文件,请确认文件系统tar是否解压.", busyboxUrl)
		}
	}
}

func createWriteLayer() {
	writeURL := constVar.WriteLayer
	os.Mkdir(writeURL, 0777)
}

func createMountPoint() error {
	logrus.Infof("开始创建挂载点目录%s.", constVar.MntURL)
	os.Mkdir(constVar.MntURL, 0777)
	cmd := exec.Command("mount", "-t", "aufs", "-o", constVar.MountAufsDirs, "none", constVar.MntURL)
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
func volumeMapping(volume string) {
	if volume != "" {
		volumeUrls := strings.Split(volume, ":")
		if len(volumeUrls) == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			mountVolume(volumeUrls)
		} else {
			logrus.Errorf("volume参数配置错误,请使用v1:v2这样的格式.")
		}
	} else {
		logrus.Info("没有挂载卷数据.")
	}
}

func mountVolume(volumeUrls []string) {
	parentUrl := volumeUrls[0]
	os.Mkdir(parentUrl, 0777)

	containerUrl := volumeUrls[1]
	containerVolumeURL := constVar.MntURL + containerUrl
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
