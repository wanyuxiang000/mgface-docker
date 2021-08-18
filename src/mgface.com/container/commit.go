package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"mgface.com/constVar"
	"os/exec"
)

func CommitContainer(containerName, imageName string) {
	mntURL := fmt.Sprintf(constVar.MntURL, containerName)
	imageTar := fmt.Sprintf(constVar.ImageStoreURL, imageName)
	logrus.Infof("提交容器:%s (容器位置:%s)转存为镜像文件:%s", containerName, mntURL, imageTar)
	//打包/root/mnt下面的所有文件
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		logrus.Errorf("打包tar文件失败:%v", err)
	}
}
