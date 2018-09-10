package aufs

import (
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

//NewWorkspace("/root","/root/mnt")
func NewWorkspace(rootURL string, mntURL string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)

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
	logrus.Infof("mount->dirs：%s",dirs)
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
