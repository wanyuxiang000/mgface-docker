package net

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
	"os"
	"path"
)

var (
	defaultNetworkPath = "/var/run/mgface-docker/network/network/"
	drivers = map[string]NetworkDriver{}
	networks = map[string]*Network{}
)
type Network struct {
	Name string //网络名称
	IpRange string //地址段
	Driver string //网络驱动名
}

type EndPoint struct {
	Id string `json:"id"`
	Device netlink.Veth `json:"dev"`
	IpAddress net.IPAddr `json:"ip"`
	macAddress net.HardwareAddr `json:"mac"`
	PortMapping []string `json:"portmapping"`
	NetWork *Network
}


type NetworkDriver interface {
	Name() string //驱动名
	Create(subnet string,name string) (*Network,error) //创建网络
	Delete(network Network) error //删除网络
	//连接容器网络端点到网络
	Connect(work *Network,point *EndPoint) error
	//从网络上移除容器端点
	Disconnet(work *Network,point *EndPoint) error
}

func CreateNetwork(driver, subnet, name string) error {
	//ParseCIDR 是 Galang net 包的函数，功能是将网段的字符串转换成 net.IPNet 的对象
	_, cidr, _ := net.ParseCIDR(subnet)
	//／／通过 IPAM 分配网关 IP，获取到网段中第一个 IP 作为网关的凹，下面几节会具体介绍它的实现
	ip, err := ipAllocator.Allocate(cidr)
	if err != nil {
		return err
	}
	cidr.IP = ip

	nw, err := drivers[driver].Create(cidr.String(), name)
	if err != nil {
		return err
	}
	//／／保存网络信息，将网络的信息保存在文件系统中，以便查询和在网络上连接网络端点
	return nw.dump(defaultNetworkPath)
}
func (nw *Network) load(dumpPath string) error {
	nwConfigFile, err := os.Open(dumpPath)
	defer nwConfigFile.Close()
	if err != nil {
		return err
	}
	nwJson := make([]byte, 2000)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		return err
	}

	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		logrus.Errorf("Error load nw info", err)
		return err
	}
	return nil
}

func (nw *Network) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		} else {
			return err
		}
	}

	nwPath := path.Join(dumpPath, nw.Name)
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		logrus.Errorf("error：", err)
		return err
	}
	defer nwFile.Close()

	nwJson, err := json.Marshal(nw)
	if err != nil {
		logrus.Errorf("error：", err)
		return err
	}

	_, err = nwFile.Write(nwJson)
	if err != nil {
		logrus.Errorf("error：", err)
		return err
	}
	return nil
}