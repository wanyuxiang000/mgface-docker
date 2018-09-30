package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mgface.com/container"
)

var CommitCommand = cli.Command{
	Name:  "commit",
	Usage: "提交一个容器成为新镜像,命令:commit [contianerName]  [imageName]",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 2 {
			return fmt.Errorf("参数必须为commit 容器名称  镜像名称")
		}
		containerName := context.Args().Get(0)
		imageName := context.Args().Get(1)
		container.CommitContainer(containerName,imageName)
		return nil
	},
}