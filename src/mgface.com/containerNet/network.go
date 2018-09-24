package containerNet

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	. "mgface.com/constVar"
	"mgface.com/containerInfo"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
)

var (
	drivers  = map[string]NetworkDriver{}
	networks = map[string]*Network{}
)

type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`
	IPAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	Network     *Network
	PortMapping []string
}

type Network struct {
	Name   string
	IpNet  *net.IPNet
	Driver string
}

func (network *Network) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		} else {
			return err
		}
	}

	nwPath := path.Join(dumpPath, network.Name)
	networkFile, _ := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	defer networkFile.Close()
	content, _ := json.MarshalIndent(network, "", "  ")
	content = append(content, []byte("\n")...)
	networkFile.Write(content)
	return nil
}

func (network *Network) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, network.Name)); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		return os.Remove(path.Join(dumpPath, network.Name))
	}
}

func (network *Network) load(dumpPath string) error {
	configFile, _ := os.Open(dumpPath)
	defer configFile.Close()
	nwJson := make([]byte, 1024*1024)
	n, _ := configFile.Read(nwJson)
	if err := json.Unmarshal(nwJson[:n], network); err != nil {
		logrus.Infof("解析文件 [%s] 失败,请核实信息:%+v.", dumpPath,err)
	}
	return nil
}

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
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
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
	}

	if err := ipAllocator.Release(network.IpNet, &network.IpNet.IP); err != nil {
		return fmt.Errorf("错误的释放Network的IP地址: %s", err)
	}

	if err := drivers[network.Driver].Delete(*network); err != nil {
		return fmt.Errorf("错误的移除网络驱动: %s", err)
	}

	//todo 同时需要移除subnet子网络
	return network.remove(DefaultNetworkPath)
}

func enterContainerNetns(enLink *netlink.Link, cinfo *containerInfo.ContainerInfo) func() {
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", cinfo.Pid), os.O_RDONLY, 0)
	if err != nil {
		logrus.Errorf("error get container net namespace, %v", err)
	}

	nsFD := f.Fd()
	runtime.LockOSThread()

	// 修改veth peer 另外一端移到容器的namespace中
	if err = netlink.LinkSetNsFd(*enLink, int(nsFD)); err != nil {
		logrus.Errorf("error set link netns , %v", err)
	}

	// 获取当前的网络namespace
	origns, err := netns.Get()
	if err != nil {
		logrus.Errorf("error get current netns, %v", err)
	}

	// 设置当前进程到新的网络namespace，并在函数执行完成之后再恢复到之前的namespace
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		logrus.Errorf("error set netns, %v", err)
	}
	return func() {
		netns.Set(origns)
		origns.Close()
		runtime.UnlockOSThread()
		f.Close()
	}
}

func configEndpointIpAddressAndRoute(ep *Endpoint, cinfo *containerInfo.ContainerInfo) error {
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)
	if err != nil {
		return fmt.Errorf("fail config endpoint: %v", err)
	}

	defer enterContainerNetns(&peerLink, cinfo)()

	interfaceIP := *ep.Network.IpNet
	interfaceIP.IP = ep.IPAddress

	if err = setInterfaceIP(ep.Device.PeerName, interfaceIP.String()); err != nil {
		return fmt.Errorf("%v,%s", ep.Network, err)
	}

	if err = setInterfaceUP(ep.Device.PeerName); err != nil {
		return err
	}

	if err = setInterfaceUP("lo"); err != nil {
		return err
	}

	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        ep.Network.IpNet.IP,
		Dst:       cidr,
	}

	if err = netlink.RouteAdd(defaultRoute); err != nil {
		return err
	}

	return nil
}

func configPortMapping(ep *Endpoint, cinfo *containerInfo.ContainerInfo) error {
	for _, pm := range ep.PortMapping {
		portMapping := strings.Split(pm, ":")
		if len(portMapping) != 2 {
			logrus.Errorf("port mapping format error, %v", pm)
			continue
		}
		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
			portMapping[0], ep.IPAddress.String(), portMapping[1])
		cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
		//err := cmd.Run()
		output, err := cmd.Output()
		if err != nil {
			logrus.Errorf("iptables Output, %v", output)
			continue
		}
	}
	return nil
}

func Connect(networkName string, cinfo *containerInfo.ContainerInfo) error {
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("No Such Network: %s", networkName)
	}

	// 分配容器IP地址
	ip, err := ipAllocator.Allocate(network.IpNet)
	if err != nil {
		return err
	}

	// 创建网络端点
	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", cinfo.Id, networkName),
		IPAddress:   ip,
		Network:     network,
		PortMapping: cinfo.PortMapping,
	}
	// 调用网络驱动挂载和配置网络端点
	if err = drivers[network.Driver].Connect(network, ep); err != nil {
		return err
	}
	// 到容器的namespace配置容器网络设备IP地址
	if err = configEndpointIpAddressAndRoute(ep, cinfo); err != nil {
		return err
	}

	return configPortMapping(ep, cinfo)
}

func Disconnect(networkName string, cinfo *containerInfo.ContainerInfo) error {
	return nil
}
