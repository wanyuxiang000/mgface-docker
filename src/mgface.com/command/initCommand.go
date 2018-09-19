package command

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mgface.com/container"
)

var InitCommand = cli.Command{
	Name:  "init",
	Usage: "在容器中做初始化功能,不要在外部进行调用这个方法",
	Action: func(context *cli.Context) error {
		logrus.Infof("初始化容器...")
		//cmd := context.Args().Get(0)
		//logrus.Infof("命令->%s", cmd)
		err := container.RunContainerInitProcess()
		return err
	},
}