package aufs

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"mgface.com/constVar"
	"os"
	"os/exec"
	"strings"
)

func DeleteFileSystem(volume string, containerName string) {
	logrus.Infof("1)删除挂载卷...")
	deleteVolumeMapping(volume, containerName)
	logrus.Infof("2)卸载挂载点...")
	deleteMountPoint(containerName)
	logrus.Infof("3)删除读写层...")
	deleteWriteLayer(containerName)
	//todo 准备清除掉tar解压的文件系统
	logrus.Info("删除只读层...")
	deleteReadOnlyLayer(containerName)

}

func deleteReadOnlyLayer(containerName string) {
	busyboxUrl := fmt.Sprintf(constVar.FileSystemURL, containerName)
	if exist, _ := pathExit(busyboxUrl); exist {
		os.RemoveAll(busyboxUrl)
	}
}

func deleteVolumeMapping(volume string, containerName string) {
	if volume != "" {
		volumeUrls := strings.Split(volume, ":")
		if len(volumeUrls) == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			deleteMountPointWithVolume(volumeUrls, containerName)
		}
	} else {
		logrus.Info("没有挂载卷数据.")
	}
}

func deleteMountPointWithVolume(volumeURL []string, containerName string) {
	containerUrl := fmt.Sprintf(constVar.MntURL, containerName) + volumeURL[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Infof("umount [%s]失败:%v", containerUrl, err)
	}
}

func deleteMountPoint(containerName string) {
	mntURL := fmt.Sprintf(constVar.MntURL, containerName)
	cmd := exec.Command("umount", mntURL)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Infof("umount发生异常:%v", err)
	}
	logrus.Infof("删除容器挂载文件系统[%v].", mntURL)
	os.RemoveAll(mntURL)
}

func deleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(constVar.WriteLayer, containerName)
	logrus.Infof("删除容器 [%v] 可写层.", writeURL)
	os.RemoveAll(writeURL)
}
