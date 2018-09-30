package containerNet

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"mgface.com/constVar"
	"net"
	"os"
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

	if err := driver.initBridge(network); err != nil {
		log.Errorf("错误的初始化bridge: %v", err)
		return network, err
	}
	return network, nil
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
	log.Info("创建Veth接口的配置.")
	la := netlink.NewLinkAttrs()
	//由于 Linux 接口名的限制,名字取endpoint ID的前5位
	la.Name = endpoint.ID[:5]
	//通过设置Veth接口的master属性,设置这个Veth的一端挂载到网络对应的Linux Bridge上
	la.MasterIndex = bridge.Attrs().Index
	//创建Veth对象,通过PeerName配置Veth另外一端的接口名
	//配置Veth另外一端的名字 cif - {endpoint ID 的前 5 位｝
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
	log.Info("启动Veth接口.")
	if err = netlink.LinkSetUp(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}
	return nil
}

func (driver *BridgeNetworkDriver) Disconnect(network Network, endpoint *Endpoint) error {
	return nil
}

func (driver *BridgeNetworkDriver) initBridge(network *Network) error {
	log.Info("1. 创建 Bridge 虚拟设备.")
	bridgeName := network.Name
	if err := createBridgeInterface(bridgeName); err != nil {
		return fmt.Errorf("错误的添加bridge： %s, Error: %v", bridgeName, err)
	}

	log.Info("2.设置Bridge设备的地址和路由.")
	ipNet := *network.IpNet
	ipNet.IP = network.IpNet.IP

	if err := setInterfaceIP(bridgeName, ipNet.String()); err != nil {
		return fmt.Errorf("错误的分配一个地址: %s 在 bridge: %s 异常信息: %v", ipNet, bridgeName, err)
	}
	log.Info("3.启动Bridge设备.")
	if err := setInterfaceUP(bridgeName); err != nil {
		return fmt.Errorf("bridge %s 设备启动发生错误: %+v", bridgeName, err)
	}

	log.Info("4.设置iptabels的SNAT规则.")
	if err := setupIPTables(bridgeName, network.IpNet); err != nil {
		return fmt.Errorf("%s 错误的设置iptables异常信息为: %v", bridgeName, err)
	}

	log.Info("5.开启linux的数据转发功能.")
	ipv4File, _ := os.OpenFile(constVar.IP4Forward, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	defer ipv4File.Close()
	ipv4File.WriteString("1")

	return nil
}

func createBridgeInterface(bridgeName string) error {

	ipface, err := net.InterfaceByName("docker0")
	if ipface!=nil {
		return errors.New("docker0设备存在,存在网络配置驱动冲突,请卸载docker/或者检查iptables策略.")
	}
	//先检查是否己经存在了这个同名的 Bridge 设备
	ipface, err = net.InterfaceByName(bridgeName)
	//如果已经存在或者报错则返回创建错误
	if ipface != nil {
		return errors.New("设备存在.创建失败.")
	}
	if err != nil && !strings.Contains(err.Error(), "no such network interface") {
		log.Errorf("创建网桥设备出错:%+v",err)
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

//设置一个网络接口的IP地址
func setInterfaceIP(name string, ipRange string) error {

	iface, _ := netlink.LinkByName(name)
	ipNet, _ := netlink.ParseIPNet(ipRange)
	addr := &netlink.Addr{IPNet: ipNet, Peer: ipNet, Label: "", Flags: 0, Scope: 0, Broadcast: nil}
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
		return fmt.Errorf("错误的启用设备 %s: %v", interfaceName, err)
	}
	return nil
}

func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	log.Infof("添加的nat映射规则:%s", iptablesCmd)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output, %v", output)
	}
	return err
}
