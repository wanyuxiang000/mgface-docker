package container

import (
	"fmt"
	"io/ioutil"
	"mgface.com/constVar"
	"mgface.com/containerInfo"
	"os"
	"text/tabwriter"
)

func ListContainers() {
	dirURL := fmt.Sprintf(constVar.DefaultInfoLocation, "")
	dirURL = dirURL[:len(dirURL)-1]
	files, _ := ioutil.ReadDir(dirURL)
	var containers []*containerInfo.ContainerInfo
	for _, file := range files {
		tmp, _ := containerInfo.GetContainerInfo(file)
		containers = append(containers, tmp)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	//打印信息
	fmt.Fprintf(w, "*****************信息展现start*****************\n")
	fmt.Fprintf(w, "ID\t容器名称\tPID\t状态\t命令\t创建时间\t结束时间\t\n")
	for _, cinfo := range containers {
		fmt.Fprintf(w,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
			cinfo.Id,
			cinfo.Name,
			cinfo.Pid,
			cinfo.Status,
			cinfo.Command,
			cinfo.CreatedTime,
			cinfo.StoppedTime)
	}
	fmt.Fprintf(w, "*****************信息展现end*****************\n")
	//刷新输出流缓冲区
	w.Flush()
}