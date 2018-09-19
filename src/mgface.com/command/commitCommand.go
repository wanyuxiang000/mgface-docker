package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mgface.com/container"
)

var CommitCommand = cli.Command{
	Name:  "commit",
	Usage: "提交一个容器成为新镜像",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("错误的容器名称")
		}
		imageName := context.Args().Get(0)
		container.CommitContainer(imageName)
		return nil
	},
}