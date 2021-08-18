package command

import (
	"fmt"
	"github.com/sirupsen/logrus"
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
	cli.StringSliceFlag{
		Name:  "p",
		Usage: "端口映射",
	},
}

var RunCommand = cli.Command{
	Name:  "run",
	Usage: "示例指令:docker run -it|-d --name mynginx -v /root/abc:/abc -net mgface0 -p 8989:8989 镜像文件系统(mginx) 后台运行命令(nginx)",
	Flags: flag,
	Action: func(ctx *cli.Context) error {

		if len(ctx.Args()) < 1 {
			return fmt.Errorf("错误的容器参数")
		}

		tty := ctx.Bool("it")

		var cmdArray []string
		for _, arg := range ctx.Args() {
			logrus.Infof("获取到参数:%s", arg)
			cmdArray = append(cmdArray, arg)
		}

		logrus.Infof("入参tty:%t,命令:%s", tty, cmdArray)
		//mgface -d --name admin -p 80:80 busybox top
		//后面的busybox top 最少需要2个
		if len(cmdArray) < 2 {
			return fmt.Errorf("错误的容器参数,需要指定引用的文件系统.")
		}

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

		//获得volume配置
		volume := ctx.String("v")
		//获得环境变量
		envs := ctx.StringSlice("e")
		//连接的网络
		network := ctx.String("net")
		//获取端口映射
		portMapping := ctx.StringSlice("p")
		container.RunContainer(tty, cmdArray, resconfig, volume, containerName, envs, network, portMapping)
		return nil
	},

	//todo 配置默认的网桥  run启动的时候默认使用
}
