package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mgface.com/containerNet"
)

var NetworkCommand = cli.Command{
	Name:      "network",
	ShortName: "net",
	Usage:     "创建网络",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "创建容器的网络",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "设置driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "设置subnet子网络",
				},
			},
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("错误的network名称")
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
			Name:        "remove",
			ShortName:   "rm",
			Usage:       "移除网络",
			Description: "查看iptables配置的MASQUERADE规则[iptables -t nat -vnL POSTROUTING]",
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("错误的network名称")
				}
				containerNet.InitNetworkAndNetdriver()
				containerNet.DeleteNetwork(context.Args().Get(0))
				return nil
			},
		},
	},
}
