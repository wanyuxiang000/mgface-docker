package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mgface.com/container"
)

var RmCommand = cli.Command{
	Name:  "rm",
	Usage: "移除容器",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("rm的参数必须不少于1")
		}
		containerName := ctx.Args().Get(0)
		return container.RemoveContainer(containerName)
	},
}
