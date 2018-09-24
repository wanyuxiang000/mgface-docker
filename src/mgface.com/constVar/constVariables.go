package constVar

const (
	Cmd = "/root/mnt/%s/" //进程切换到子进程的文件目录
	MntURL = "/root/mnt/%s/" //创建挂载点
	FileSystemTarURL = "/root/busybox.tar" //使用的文件系统的tar目录，解压到FileSystemURL指定的目录
	FileSystemURL = "/root/busybox/%s/" //使用的文件系统存放的路径
	WriteLayer = "/root/writeLayer/%s/" //创建可写层
	MountAufsDirs = "dirs=" + WriteLayer + ":" + FileSystemURL //挂载的Aufs文件系统
	CgroupName = "mgfaceCgroup/%s"
	ImageStoreURL = "/root/%s.tar"
	//容器属性
	RUNNING             = "running"
	STOP                = "stopped"
	DefaultInfoLocation = "/var/run/mgface-docker/%s/"
	ConfigName          = "config.json"
	ContainerLog        = "container.log"
	//网络
    IpamDefaultAllocatorPath = "/var/run/mgface-docker-net/ipam/subnet.json"
    //网络路径的存放位置
	DefaultNetworkPath = "/var/run/mgface-docker-net/network/"
)