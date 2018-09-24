package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mgface.com/container"
)

var StopCommand = cli.Command{
	Name:  "stop",
	Usage: "停止容器.",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("参数错误，请使用stop containerName")
		}
		containerName := ctx.Args().Get(0)
		return container.StopContainer(containerName)
	},
}
