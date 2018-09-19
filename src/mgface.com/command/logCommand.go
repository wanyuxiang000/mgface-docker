package command

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"io/ioutil"
	"mgface.com/containerInfo"
	"os"
)

var LogCommand = cli.Command{
	Name:  "logs",
	Usage: "查看指定容器的日志文件",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("请指定容器的名称.")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

func logContainer(containerName string) {
	dirURL := fmt.Sprintf(containerInfo.DefaultInfoLocation, containerName)
	logFile := dirURL + containerInfo.ContainerLog
	file, err := os.Open(logFile)
	if err != nil {
		logrus.Errorf("错误的读取文件%s,发生的异常为:%v", file, err)
	}
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	fmt.Fprintf(os.Stdout, string(content))
}