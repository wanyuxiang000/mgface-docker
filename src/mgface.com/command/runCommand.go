package command

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mgface.com/cgroup"
	"mgface.com/container"
	"os"
)

var flag = []cli.Flag{
	cli.BoolFlag{
		Name:  "it",
		Usage: "启用tty",
	},
	cli.StringFlag{
		Name:  "m",
		Usage: "内存限制",
	},
	cli.StringFlag{
		Name:  "cpushare",
		Usage: "cpushare限制",
	},
	cli.StringFlag{
		Name:  "cpuset",
		Usage: "cpuset限制",
	},
	cli.StringFlag{
		Name:  "v",
		Usage: "volume,挂载多个文件请使用\",\"隔开",
	},
	cli.BoolFlag{
		Name:  "d",
		Usage: "detach container",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "指定容器名称.",
	},
	cli.StringSliceFlag{
		Name:  "e",
		Usage: "设置环境变量.",
	},
	cli.StringFlag{
		Name:  "net",
		Usage: "指定连接的网络",
	},
	cli.StringFlag{
		Name:  "p",
		Usage: "端口映射",
	},
}

var RunCommand = cli.Command{
	Name:  "run",
	Usage: "创建一个容器使用cgroup和namespace,指令为docker run -it [command]",
	Flags: flag,
	Action: func(ctx *cli.Context) error {

		if len(ctx.Args()) < 1 {
			return fmt.Errorf("错误的容器参数")
		}
		var cmdArray []string
		for _, arg := range ctx.Args() {
			logrus.Infof("获取到参数:%s", arg)
			cmdArray = append(cmdArray, arg)
		}

		tty := ctx.Bool("it")

		detach := ctx.Bool("d")

		containerName := ctx.String("name")

		if detach && tty {
			logrus.Errorf("-it 和 -d 不能同时存在.")
			os.Exit(-1)
		}

		resconfig := &cgroup.ResouceConfig{
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
			MemoryLimit: ctx.String("m"),
		}
		logrus.Infof("入参tty:%t,命令:%s", tty, cmdArray)
		//获得volume配置
		volume := ctx.String("v")

		//获得环境变量
		envs := ctx.StringSlice("e")

		//连接的网络
		network:=ctx.String("net")
		//获取端口映射
		portMapping :=ctx.StringSlice("p")
		container.RunContainer(tty, cmdArray, resconfig, volume, containerName, envs,network,portMapping)
		return nil
	},
}
