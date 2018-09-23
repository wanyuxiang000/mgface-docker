package container

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"mgface.com/constVar"
	"os"
)

func LogContainer(containerName string) {
	dirURL := fmt.Sprintf(constVar.DefaultInfoLocation, containerName)
	logFile := dirURL + constVar.ContainerLog
	file, err := os.Open(logFile)
	if err != nil {
		logrus.Errorf("错误的读取文件%s,发生的异常为:%v", file, err)
	}
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	fmt.Fprintf(os.Stdout, string(content))
}