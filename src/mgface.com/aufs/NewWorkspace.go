package aufs

import (
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

//NewWorkspace("/root","/root/mnt")
func NewWorkspace(rootURL string, mntURL string, volume string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
	if volume != "" {
		volumeUrls := strings.Split(volume, ":")
		if len(volumeUrls) == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			MountVolume(rootURL, mntURL, volumeUrls)
		} else {
			logrus.Errorf("volume参数配置错误，请使用master:slave这样的格式.")
		}
	}

}

func MountVolume(rootURL, mntURL string, volumeUrls []string) {
	parentUrl := volumeUrls[0]
	os.Mkdir(parentUrl, 0777)

	containerUrl := volumeUrls[1]
	containerVolumeURL := mntURL + containerUrl
	os.Mkdir(containerVolumeURL, 0777)
	//把宿主机文件挂到容器挂载点上
	dirs := "dirs=" + parentUrl
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount volume 错误.%v", err)
	}
}

func CreateReadOnlyLayer(rootURL string) {
	busyboxUrl := rootURL + "/busybox/"
	busyboxTarURL := rootURL + "/busybox.tar"
	exist, err := PathExit(busyboxUrl)
	if err != nil {
		logrus.Infof("错误的目录:%s,发生异常:%v", busyboxUrl, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxUrl, 0777); err != nil {
			logrus.Errorf("创建目录%s发生异常%v", busyboxUrl, err)
		}
		logrus.Infof("目录:%s,tar包:%s", busyboxUrl, busyboxTarURL)
		exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxUrl).CombinedOutput()
	}
}

func CreateWriteLayer(rootRUL string) {
	writeURL := rootRUL + "/writeLayer/"
	os.Mkdir(writeURL, 0777)
}

func CreateMountPoint(rootURL string, mntURL string) {
	os.Mkdir(mntURL, 0777)
	dirs := "dirs=" + rootURL + "/writeLayer:" + rootURL + "/busybox"
	logrus.Infof("mount->dirs：%s", dirs)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("发错致命错误:%v", err)
	}
}

func PathExit(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}

func DeleteWorkSpace(rootURL string, mntURL string, volume string) {
	if volume != "" {
		volumeUrl := strings.Split(volume, ",")
		if len(volumeUrl) == 2 && volumeUrl[0] != "" && volumeUrl[1] != "" {
			DeleteMountPointWithVolume(rootURL, mntURL, volumeUrl)
		}
	} else {
		DeleteMountPoint(rootURL, mntURL)
	}
	DeleteWriteLayer(rootURL)

}

func DeleteMountPointWithVolume(rootURL string, mntURL string, volumeURL []string) {
	containerUrl := mntURL + volumeURL[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	//卸载整个容器系统的挂载点
	cmd = exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	logrus.Infof("删除容器文件系统的挂载点...")
	os.RemoveAll(mntURL)
}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	os.RemoveAll(rootURL + "/busybox")
	os.RemoveAll(mntURL)
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "/writeLayer/"
	os.RemoveAll(writeURL)
}
