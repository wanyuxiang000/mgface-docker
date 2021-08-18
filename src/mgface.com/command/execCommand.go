package command

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mgface.com/container"
	"os"
)

var ExecCommand = cli.Command{
	Name:  "exec",
	Usage: "进入容器",
	Action: func(context *cli.Context) error {
		if os.Getenv(container.EnvExecPid) != "" {
			logrus.Infof("pid(%d)进行了自身回调.", os.Getpid())
			return nil
		}
		if len(context.Args()) < 2 {
			return fmt.Errorf("参数错误.")
		}
		containerName := context.Args().Get(0)
		var commandArray []string
		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		//进入容器
		container.ExecContainer(containerName, commandArray)
		return nil
	},
}
