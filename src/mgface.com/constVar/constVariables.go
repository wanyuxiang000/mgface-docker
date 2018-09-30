package constVar

const (
	Cmd = "/root/mnt/%s/" //进程切换到子进程的文件目录
	MntURL = "/root/mnt/%s/" //创建挂载点

	//提示:文件系统可以从docker里面下载，docker export -o xxx.tar 容器ID
	//例如/root/busybox（文件系统）.tar
	FileSystemTarURL = "/root/%s.tar" //使用的文件系统的tar目录，解压到FileSystemURL指定的目录
	//例如/root/busybox(文件系统)/admin（容器）
	FileSystemURL = "/root/%s/%s/" //使用的文件系统存放的路径
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
	//IP4开启转发
	IP4Forward = "/proc/sys/net/ipv4/ip_forward"


	//网桥配置
	BridgeType = "bridge"
	Subnet     = "172.18.0.0/24"
	BridgeName = "mgface0"
)