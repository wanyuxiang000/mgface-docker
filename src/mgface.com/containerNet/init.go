package containerNet

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	. "mgface.com/constVar"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

func InitNetworkAndNetdriver() error {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver

	if _, err := os.Stat(DefaultNetworkPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(DefaultNetworkPath, 0644)
		} else {
			return err
		}
	}

	filepath.Walk(DefaultNetworkPath, func(networkPath string, info os.FileInfo, err error) error {
		logrus.Infof("读取到文件:%s", networkPath)
		if strings.HasSuffix(networkPath, "/") {
			return nil
		}
		_, networkName := path.Split(networkPath)
		network := &Network{
			Name: networkName,
		}
		network.load(networkPath)
		networks[networkName] = network
		return nil
	})
	return nil
}

func CreateNetwork(driver, subnet, name string) error {
	_, ipNet, _ := net.ParseCIDR(subnet)
	ip, _ := ipAllocator.Allocate(ipNet)
	ipNet.IP = ip

	network, _ := drivers[driver].Create(ipNet.String(), name)
	return network.dump(DefaultNetworkPath)
}

func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 3, 3, ' ', 0)
	fmt.Fprint(w, "网络名称\tIP网络\t网络驱动\n")
	for _, network := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			network.Name,
			network.IpNet.String(),
			network.Driver,
		)
	}
	w.Flush()
}

func DeleteNetwork(networkName string) error {
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("没有匹配到Network: %s", networkName)
	}else {
		//删除子网络
		delete(networks,networkName)
	}

	if err := ipAllocator.Release(network.IpNet, &network.IpNet.IP); err != nil {
		return fmt.Errorf("错误的释放Network的IP地址: %s", err)
	}

	if err := drivers[network.Driver].Delete(*network); err != nil {
		return fmt.Errorf("错误的移除网络驱动: %s", err)
	}
	

	//删除ipam数据
	_, subnet, _ := net.ParseCIDR(network.IpNet.String())
	delete(*ipAllocator.Subnets,subnet.String())
	ipAllocator.dump()
	//移除subnet子网络文件
	return network.remove(DefaultNetworkPath)
}
