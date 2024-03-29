package command

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mgface.com/container"
)

var InitCommand = cli.Command{
	Name:  "init",
	Usage: "在容器中做初始化功能,不要在外部进行调用这个方法",
	Action: func(context *cli.Context) error {
		logrus.Infof("初始化容器...")
		err := container.RunContainerInitProcess()
		return err
	},
}
