package container

import (
	"github.com/Sirupsen/logrus"
	"os/exec"
)

func CommitContainer(imageName string) {
	mntURL := "/root/mnt"
	imageTar := "/root/" + imageName + ".tar"
	logrus.Infof("镜像文件:%s", imageTar)
	//打包/root/mnt下面的所有文件
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		logrus.Errorf("打包tar文件失败:%v", err)
	}
}
