package constVar

const (
	Cmd = "/root/mnt" //进程切换到子进程的文件目录
	MntURL = "/root/mnt" //创建挂载点
	FileSystemTarURL = "/root/busybox.tar" //使用的文件系统的tar目录，解压到FileSystemURL指定的目录
	FileSystemURL = "/root/busybox/" //使用的文件系统存放的路径
	WriteLayer = "/root/writeLayer/" //创建可写层
	MountAufsDirs = "dirs=" + WriteLayer + ":" + FileSystemURL //挂载的Aufs文件系统
	CgroupName = "mgfaceCgroup"
)