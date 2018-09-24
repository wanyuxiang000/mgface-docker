package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mgface.com/containerNet"
)

var NetworkCommand = cli.Command{
	Name:  "network",
	Usage: "创建网络",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "创建容器的网络",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet 子网络",
				},
			},
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				containerNet.InitNetworkAndNetdriver()
				containerNet.CreateNetwork(context.String("driver"), context.String("subnet"), context.Args().Get(0))
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "显示容器网络",
			Action: func(context *cli.Context) error {
				containerNet.InitNetworkAndNetdriver()
				containerNet.ListNetwork()
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove container network",
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				containerNet.InitNetworkAndNetdriver()
				err := containerNet.DeleteNetwork(context.Args()[0])
				if err != nil {
					return fmt.Errorf("remove network error: %+v", err)
				}
				return nil
			},
		},
	},
}
