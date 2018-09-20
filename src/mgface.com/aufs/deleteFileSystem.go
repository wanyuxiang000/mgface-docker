package aufs

import (
	"github.com/Sirupsen/logrus"
	"mgface.com/constVar"
	"os"
	"os/exec"
	"strings"
)

func DeleteFileSystem(volume string) {
	logrus.Infof("1)卸载挂载点...")
	deleteMountPoint()
	logrus.Infof("2)删除读写层...")
	deleteWriteLayer()
	logrus.Infof("3)删除挂载卷...")
	deleteVolumeMapping(volume)
}

func deleteVolumeMapping(volume string) {
	if volume != "" {
		volumeUrls := strings.Split(volume, ":")
		if len(volumeUrls) == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			deleteMountPointWithVolume(volumeUrls)
		}
	} else {
		logrus.Info("没有挂载卷数据.")
	}
}

func deleteMountPointWithVolume(volumeURL []string) {
	containerUrl := constVar.MntURL + volumeURL[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Infof("umount [%s]失败:%v", containerUrl, err)
	}
}

func deleteMountPoint() {
	mntURL := constVar.MntURL
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

func deleteWriteLayer() {
	writeURL := constVar.WriteLayer
	logrus.Infof("删除容器%v可写层.", writeURL)
	os.RemoveAll(writeURL)
}
