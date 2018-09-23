package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mgface.com/command"
	"os"
)

const usage = `mgface是一个简单的容器应用,我们的目的是怎么把玩docker?
 *      ┌─┐       ┌─┐
 *   ┌──┘ ┴───────┘ ┴──┐
 *   │                 │
 *   │       ───       │
 *   │  ─┬┘       └┬─  │
 *   │                 │
 *   │       ─┴─       │
 *   │                 │
 *   └───┐         ┌───┘
 *       │         │
 *       │         │
 *       │         │
 *       │         └──────────────┐
 *       │                        │
 *       │                        ├─┐
 *       │                        ┌─┘
 *       │                        │
 *       └─┐  ┐  ┌───────┬──┐  ┌──┘
 *         │ ─┤ ─┤       │ ─┤ ─┤
 *         └──┴──┘       └──┴──┘
 *                神兽保佑
 *               代码无BUG!
`

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Sprintf("发生致命错误.%v\n", e)
		}
	}()
	app := cli.NewApp()
	app.Name = "mgface"
	app.Version = "1.0.0"
	app.Author = "Yuxiang Wan"
	app.Copyright = "mgface@2018-∞"
	app.Usage = usage
	app.Email = "15622535353@163.com"
	app.Commands = []cli.Command{
		command.RunCommand,
		command.InitCommand,
		command.CommitCommand,
		command.ListCommand,
		command.LogCommand,
		command.ExecCommand,
		command.StopCommand,
		command.RmCommand,
		command.StartCommand,
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
