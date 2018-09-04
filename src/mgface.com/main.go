package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mgface.com/container"
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

var runCommand = cli.Command{
	Name:  "run",
	Usage: "创建一个容器使用cgroup和namespace,指令为docker run -it [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "启用tty",
		},
	},
	Action: func(ctx *cli.Context) error {
		logrus.Infof("获取到参数:%s", ctx.Args)
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("错误的容器参数")
		}
		cmd := ctx.Args().Get(0)
		tty := ctx.Bool("it")
		container.Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "在容器中做初始化功能,不要在外部进行调用这个方法",
	Action: func(context *cli.Context) error {
		logrus.Infof("初始化容器...")
		cmd := context.Args().Get(0)
		logrus.Infof("命令->%s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
