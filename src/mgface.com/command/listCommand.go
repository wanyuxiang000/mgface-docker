package command

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"mgface.com/containerInfo"
	"os"
	"text/tabwriter"
)

var ListCommand = cli.Command{
	Name:  "ps",
	Usage: "显示所有的容器",
	Action: func(context *cli.Context) error {
		listContainers()
		return nil
	},
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