package containerNet

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
	"os/exec"
	"strings"
)

type BridgeNetworkDriver struct {
}

func (driver *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (driver *BridgeNetworkDriver) Create(subnet string, name string) (*Network, error) {
	//通过net包中的net.ParseCIDR方法,取到网段的字符串中的网关IP地址和网络IP段
	ip, ipNet, _ := net.ParseCIDR(subnet)
	ipNet.IP = ip
	network := &Network{
		Name:   name,
		IpNet:  ipNet,
		Driver: driver.Name(),
	}
	err := driver.initBridge(network)
	if err != nil {
		log.Errorf("错误的初始化bridge: %v", err)
	}
	return network, err
}

func (driver *BridgeNetworkDriver) Delete(network Network) error {
	bridgeName := network.Name
	bridge, _ := netlink.LinkByName(bridgeName)
	return netlink.LinkDel(bridge)
}

func (driver *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	bridgeName := network.Name
	bridge, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}
	//创建 Veth 接口的配置
	la := netlink.NewLinkAttrs()
	//由于 Linux 接口名的限制，名字取 endpoint ID 的前 5 位
	la.Name = endpoint.ID[:5]
	//通过设置Veth接口的master属性，设置这个Veth的一端挂载到网络对应的Linux Bridge上
	la.MasterIndex = bridge.Attrs().Index
	//创建 Veth 对象，通过 PeerName 配置 Veth 另外一端的接口名
	//配置 Veth 另外一端的名字 cif - {endpoint ID 的前 5 位｝
	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  "cif-" + la.Name,
	}
	//调用 netlink 的 LinkAdd 方法创建出这个 Veth 接口
	//因为上面指定了link的MasterIndex是网络对应的Linux Bridge
	//所以 Veth 的一端就己经挂载到了网络对应的 Linux Bridge 上
	if err = netlink.LinkAdd(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}
	//调用 netlink 的 LinkSetUp 方法，设置 Veth 启动
	//相当于 ip link set xxx up 命令
	if err = netlink.LinkSetUp(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}
	return nil
}

func (driver *BridgeNetworkDriver) Disconnect(network Network, endpoint *Endpoint) error {
	return nil
}

func (driver *BridgeNetworkDriver) initBridge(network *Network) error {
	// 1. 创建 Bridge 虚拟设备
	bridgeName := network.Name
	if err := createBridgeInterface(bridgeName); err != nil {
		return fmt.Errorf("Error add bridge： %s, Error: %v", bridgeName, err)
	}

	//2.设置Bridge设备的地址和路由
	gatewayIP := *network.IpNet
	gatewayIP.IP = network.IpNet.IP

	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		return fmt.Errorf("Error assigning address: %s on bridge: %s with an error of: %v", gatewayIP, bridgeName, err)
	}
	//3.启动Bridge设备
	if err := setInterfaceUP(bridgeName); err != nil {
		return fmt.Errorf("Error set bridge up: %s, Error: %v", bridgeName, err)
	}

	//4.设置iptabels的SNAT规则
	if err := setupIPTables(bridgeName, network.IpNet); err != nil {
		return fmt.Errorf("Error setting iptables for %s: %v", bridgeName, err)
	}

	return nil
}

func createBridgeInterface(bridgeName string) error {
	//先检查是否己经存在了这个同名的 Bridge 设备
	_, err := net.InterfaceByName(bridgeName)
	//如果已经存在或者报错则返回创建错误
	if err == nil || !strings.Contains(err.Error(), "no such network interface") {
		return err
	}

	//初始化一个netlink的Link基础对象,Link的名字即Bridge虚拟设备的名字
	link := netlink.NewLinkAttrs()
	link.Name = bridgeName
	//使用刚才创建的Link的属性创建netlink的Bridge对象
	bridge := &netlink.Bridge{LinkAttrs: link}
	//调用netlink的Linkadd方法,创建Bridge虚拟网络设备
	//netLink的Linkadd方法是用来创建虚拟网络设备的,相当于ip link add xxxx
	if err := netlink.LinkAdd(bridge); err != nil {
		return fmt.Errorf("Bridge %s 创建失败: %v", bridgeName, err)
	}
	return nil
}

//设置一个网络接口的IP地址,例如setInterfaceIP("testbridge","192.168.0.1/24")
func setInterfaceIP(bridgeName string, rawIP string) error {
	iface, _ := netlink.LinkByName(bridgeName)
	/*
		 由于 netlink.ParseIPNet是对net.ParseCIDR的一个封装,因此可以将 net.ParseCIDR的返回值的IP和net整合
		返回值中的ipNet既包含了网段的信息,192.168.0.0/24,也包含了原始的ip 192.168.0.1
	*/
	ipNet, _ := netlink.ParseIPNet(rawIP)
	addr := &netlink.Addr{IPNet: ipNet}
	//调用netlink的AddrAdd方法,配置Linux Bridge的地址和路由表。
	return netlink.AddrAdd(iface, addr)
}

// deleteBridge deletes the bridge
func (driver *BridgeNetworkDriver) deleteBridge(n *Network) error {
	bridgeName := n.Name

	// get the link
	l, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("Getting link with name %s failed: %v", bridgeName, err)
	}

	// delete the link
	if err := netlink.LinkDel(l); err != nil {
		return fmt.Errorf("Failed to remove bridge interface %s delete: %v", bridgeName, err)
	}

	return nil
}

func setInterfaceUP(interfaceName string) error {
	iface, _ := netlink.LinkByName(interfaceName)

	if err := netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("Error enabling interface for %s: %v", interfaceName, err)
	}
	return nil
}

func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output, %v", output)
	}
	return err
}
