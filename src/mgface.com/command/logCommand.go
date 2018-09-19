package command

import (
	"fmt"
	"github.com/urfave/cli"
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