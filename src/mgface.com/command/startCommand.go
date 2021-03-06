package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mgface.com/container"
)

var StartCommand = cli.Command{
	Name:  "start",
	Usage: "开始一个存在的容器.",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("参数错误，请使用start containerName")
		}
		containerName := ctx.Args().Get(0)
		return container.StartContainer(containerName)
	},
}
