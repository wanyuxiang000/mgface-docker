package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mgface.com/container"
	"mgface.com/subsystem"
	"os"
)

const usage = "mgface是一个简单的容器应用，我们的目的是搞清楚docker到底是怎么玩的？let go it"

func main() {
	app := cli.NewApp()
	app.Name = "mgface"
	app.Usage = usage
	app.Commands = []cli.Command{
		runCommand,
		initCommand,
		commitCommand,
	}
	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

var commitCommand = cli.Command{
	Name:"commit",
	Usage:"提交一个容器成为新镜像",
	Action: func(context *cli.Context) error{
		if len(context.Args())<1 {
			return fmt.Errorf("错误的容器名称")
		}
		imageName:=context.Args().Get(0)
		container.CommitContainer(imageName)
		return nil
	},
}
var runCommand = cli.Command{
	Name:  "run",
	Usage: "创建一个容器使用cgroup和namespace,指令为docker run -it [command]",
	Flags: []cli.Flag{
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
			Name:"v",
			Usage:"volume",
		},
		cli.BoolFlag{
			Name:"d",
			Usage:"detach container",
		},
	},
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

		detach:=ctx.Bool("d")

		if detach && tty {
			logrus.Errorf("-it 和 -d 不能同时存在.")
			os.Exit(-1)
		}

		resconfig := &subsystem.ResouceConfig{
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
			MemoryLimit: ctx.String("m"),
		}
		logrus.Infof("入参:%t,命令:%s", tty, cmdArray)
		//获得volume配置
		volume:=ctx.String("v")
		container.Run(tty, cmdArray, resconfig,volume)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "在容器中做初始化功能,不要在外部进行调用这个方法",
	Action: func(context *cli.Context) error {
		logrus.Infof("初始化容器...")
		//cmd := context.Args().Get(0)
		//logrus.Infof("命令->%s", cmd)
		err := container.RunContainerInitProcess()
		return err
	},
}
