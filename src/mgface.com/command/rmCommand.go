package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mgface.com/container"
)

var RmCommand = cli.Command{
	Name:  "rm",
	Usage: "mainfunc -rm 容名称",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args())<1 {
			fmt.Errorf("rm的参数必须不少于1")
			return nil
		}
		containerName:=ctx.Args().Get(0)

		//todo rm的时候还要删除挂载的文件。
		return container.RemoveContainer(containerName)
	},
}
