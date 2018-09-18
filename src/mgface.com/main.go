package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"io/ioutil"
	"mgface.com/container"
	"mgface.com/containerInfo"
	"mgface.com/subsystem"
	"os"
	"text/tabwriter"
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
		listCommand,
		logCommand,
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

var logCommand = cli.Command{
	Name:  "logs",
	Usage: "查看指定容器的日志文件",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("请指定容器的名称.")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

func logContainer(containerName string) {
	dirURL := fmt.Sprintf(containerInfo.DefaultInfoLocation, containerName)
	logFile := dirURL + containerInfo.ContainerLog
	file, err := os.Open(logFile)
	if err != nil {
		logrus.Errorf("错误的读取文件%s,发生的异常为:%v", file, err)
	}
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	fmt.Fprintf(os.Stdout, string(content))
}

func listContainers() {
	dirURL := fmt.Sprintf(containerInfo.DefaultInfoLocation, "")
	dirURL = dirURL[:len(dirURL)-1]
	files, _ := ioutil.ReadDir(dirURL)
	var containers []*containerInfo.ContainerInfo
	for _, file := range files {
		tmp, _ := containerInfo.GetContainerInfo(file)
		containers = append(containers, tmp)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	//打印信息
	fmt.Fprintf(w, "ID\t容器名称\tPID\t状态\t命令\t创建时间\t\n")
	for _, cinfo := range containers {
		fmt.Fprintf(w,
			"%s\t%s\t%s\t%s\t%s\t%s\t\n",
			cinfo.Id,
			cinfo.Name,
			cinfo.Pid,
			cinfo.Status,
			cinfo.Command,
			cinfo.CreatedTime)
	}
	fmt.Fprintf(w, "<======信息展现结束======>\n")
	//刷新输出流缓冲区
	w.Flush()
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "显示所有的容器",
	Action: func(context *cli.Context) error {
		listContainers()
		return nil
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "提交一个容器成为新镜像",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("错误的容器名称")
		}
		imageName := context.Args().Get(0)
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
			Name:  "v",
			Usage: "volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "指定容器名称.",
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

		detach := ctx.Bool("d")

		containerName := ctx.String("name")

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
		volume := ctx.String("v")
		container.Run(tty, cmdArray, resconfig, volume, containerName)
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
