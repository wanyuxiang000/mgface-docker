package containerNet

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"mgface.com/containerInfo"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)
//将容器的网络端点加入到容器的网络空间中
//并锁定当前程序所执行的线程，使当前线程进入到容器的网络空间
// 返回值是一个函数指针，执行这个返回函数才会退出容器的网络空间，回归到宿主机的网络空间
func enterContainerNetns(enLink *netlink.Link, cinfo *containerInfo.ContainerInfo) func() {
	//找到容器的 Net Namespace
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", cinfo.Pid), os.O_RDONLY, 0)
	if err != nil {
		logrus.Errorf("error get container net namespace, %v", err)
	}
	//取到文件的文件描述符
	nsFD := f.Fd()
	//锁定当前程序所执行的线程，如果不锁定操作系统线程的话
	// Go 语言的 goroutine 可能会被调度到别的线程上去
	//就不能保证一直在所需要的网络空间中了
	//所以调用runtime.LockOSThread 时要先锁定当前程序执行的线程
	runtime.LockOSThread()

	// 修改veth peer 另外一端移到容器的namespace中
	if err = netlink.LinkSetNsFd(*enLink, int(nsFD)); err != nil {
		logrus.Errorf("error set link netns , %v", err)
	}

	// 获取当前的网络namespace
	//以便后面从容器的 Net Namespace 中退出，回到原本网络的 Net Namespace 中
	origns, err := netns.Get()
	if err != nil {
		logrus.Errorf("error get current netns, %v", err)
	}

	// 调用 netns.Set 方法，将当前进程加入容器的 Net Namespace
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		logrus.Errorf("error set netns, %v", err)
	}
	//返回之前 Net Namespace 的函数
	//在容苦苦的网络空间中，执行完容器配置之后调用此函数就可以将程序恢复到原生的 Net Namespace
	return func() {
		//恢复到上面获取到的之前的 Net Namespace
		netns.Set(origns)
		//关闭 Name space 文件
		origns.Close()
		//取消对当附程序的线程锁定
		runtime.UnlockOSThread()
		//关闭 Namespace 文件
		f.Close()
	}
}
//配置容器网络端点的地址和路由
func configEndpointIpAddressAndRoute(ep *Endpoint, cinfo *containerInfo.ContainerInfo) error {
	//通过网络端点中Veth的另一端
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)
	if err != nil {
		return fmt.Errorf("fail config endpoint: %v", err)
	}
	//将容器的网络端点加入到容器的网络空间中
	//并使这个函数下面的操作都在这个网络空间中进行
	//执行完函数后，恢复为默认的网络空间
	defer enterContainerNetns(&peerLink, cinfo)()
	//获取到容器的 IP 地址及网段 ， 用于配置容器内部接口地址
	//比如容器 IP 是 192.168.1.2 ，而网络的网段是192.168.1.0/24
	//那么这里产出的 IP 字符串就是 192.168.1.2/24,用于容器内 Veth 端点配置
	interfaceIP := *ep.Network.IpNet
	interfaceIP.IP = ep.IPAddress
	//调用 setinterfaceIP 函数设置容器内 Veth 端点的 IP
	if err = setInterfaceIP(ep.Device.PeerName, interfaceIP.String()); err != nil {
		return fmt.Errorf("%v,%s", ep.Network, err)
	}
	//启动容苦苦内的 Veth 端点
	if err = setInterfaceUP(ep.Device.PeerName); err != nil {
		return err
	}
	//Net Namespace 中默认本地地址127.0.0.1的"lo"网卡是关闭状态的
	//启动它以保证容器访问自己的请求
	if err = setInterfaceUP("lo"); err != nil {
		return err
	}
	//设置容器内的外部请求都通过容器内的 Veth 端点访问
	//II 0.0.0.0/0的网段，表示所有的 IP 地址段
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")
	//构建要添加的路由数据，包括网络设备、网关 IP 及目的网段
	//相当于 route add -net 0.0.0.0/0 gw {Bridge 网桥地址｝ dev ｛容器内的 Veth 端点设备｝
	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        ep.Network.IpNet.IP,
		Dst:       cidr,
	}
	//调用 netlink 的 RouteAdd ， 添加路由到容器的网络空间
	//RouteAdd 函数相当于 route add 命令
	if err = netlink.RouteAdd(defaultRoute); err != nil {
		return err
	}

	return nil
}
//配置端口映射
func configPortMapping(ep *Endpoint, cinfo *containerInfo.ContainerInfo) error {
	//遍历容器端口映射列表
	for _, pm := range ep.PortMapping {
		//分割成宿主机的端口和容器的端口
		portMapping := strings.Split(pm, ":")
		if len(portMapping) != 2 {
			logrus.Errorf("port mapping format error, %v", pm)
			continue
		}
		//由于 iptables 没有 Go 语言版本的实现，所以采用 exec.Command 的方式直接调用命令配置
		//在 iptables 的 PREROUTING 中添加 DNAT 规则
		//将宿主机的端口请求转发到容器的地址和端口上
		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
			portMapping[0], ep.IPAddress.String(), portMapping[1])
		//执行 iptables 命令 ， 添加端口映射转发规则
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
		return fmt.Errorf("没有找到匹配的Network: %s", networkName)
	}

	// 分配容器IP地址
	ip, err := ipAddressManage.Allocate(network.IpNet)
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
	//配置端口映射信息,例如docker run -p 8080:80
	return configPortMapping(ep, cinfo)
}

func Disconnect(networkName string, cinfo *containerInfo.ContainerInfo) error {
	return nil
}
