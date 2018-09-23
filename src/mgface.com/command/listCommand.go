package command

import (
	"github.com/urfave/cli"
	"mgface.com/container"
)

var ListCommand = cli.Command{
	Name:  "ps",
	Usage: "显示所有的容器",
	Action: func(context *cli.Context) error {
		container.ListContainers()
		return nil
	},
}

