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

//将配置的容器网络端点加入到容器的网络空间中
//并锁定当前程序所执行的线程,使当前线程进入到容器的网络空间
// 返回值是一个函数指针,执行这个返回函数才会退出容器的网络空间,回归到宿主机的网络空间
func enterContainerNetns(link *netlink.Link, containerInfo *containerInfo.ContainerInfo) func() {
	//找到容器的Net Namespace
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", containerInfo.Pid), os.O_RDONLY, 0)
	if err != nil {
		logrus.Errorf("错误的获取到容器的net namespace,错误信息:%v", err)
	}
	//取到文件的文件描述符
	nsFD := f.Fd()
	//锁定当前程序所执行的线程，如果不锁定操作系统线程的话,Go语言的goroutine可能会被调度到别的线程上去
	//就不能保证一直在所需要的网络空间中了,所以调用runtime.LockOSThread 时要先锁定当前程序执行的线程
	runtime.LockOSThread()

	// 修改veth peer另外一端移到容器的namespace中
	if err = netlink.LinkSetNsFd(*link, int(nsFD)); err != nil {
		logrus.Errorf("错误设置link netns,错误信息: %v", err)
	}

	//获取当前网络的namespace
	//以便后面从容器的Net Namespace中退出,回到原本网络的Net Namespace中
	origns, err := netns.Get()
	if err != nil {
		logrus.Errorf("获取当前网络的netns发生异常:%v", err)
	}

	//设置当前进程到新的网络namespace,并在函数执行完成之后再恢复到之前的namespace
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		logrus.Errorf("错误设置netns,错误信息: %v", err)
	}
	//返回之前Net Namespace的函数
	//在容器的网络空间中,执行完容器配置之后调用此函数就可以将程序恢复到原生的 Net Namespace
	return func() {
		//恢复到上面获取到的之前的Net Namespace
		netns.Set(origns)
		//关闭Namespace 文件
		origns.Close()
		//取消对当前程序的线程锁定
		runtime.UnlockOSThread()
		//关闭Namespace 文件
		f.Close()
	}
}

//配置容器网络端点的地址和路由
func configEndpointIpAddressAndRoute(endpoint *Endpoint, containerInfo *containerInfo.ContainerInfo) error {
	//获得网络端点中Veth的另一端
	peerLink, err := netlink.LinkByName(endpoint.Device.PeerName)
	if err != nil {
		return fmt.Errorf("错误的配置endpoint: %v", err)
	}
	//将容器的网络端点加入到容器的网络空间中
	//并使这个函数下面的操作都在这个网络空间中进行,执行完函数后,恢复为默认的网络空间
	defer enterContainerNetns(&peerLink, containerInfo)()
	//获取到容器的IP地址及网段,用于配置容器内部接口地址
	//比如容器IP是192.168.1.2,而网络的网段是192.168.1.0/24
	//那么这里产出的IP字符串就是192.168.1.2/24,用于容器内Veth端点配置
	interfaceIP := *endpoint.Network.IpNet
	interfaceIP.IP = endpoint.IPAddress
	//设置容器内Veth端点的IP
	if err = setInterfaceIP(endpoint.Device.PeerName, interfaceIP.String()); err != nil {
		return fmt.Errorf("%v,%s", endpoint.Network, err)
	}
	//启动容器内的Veth端点
	if err = setInterfaceUP(endpoint.Device.PeerName); err != nil {
		return err
	}
	//Net Namespace中默认本地地址127.0.0.1的"lo"网卡是关闭状态的
	//启动它以保证容器访问自己的请求
	if err = setInterfaceUP("lo"); err != nil {
		return err
	}
	//设置容器内的外部请求都通过容器内的Veth端点访问
	//0.0.0.0/0的网段,表示所有的IP地址段
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")
	//构建要添加的路由数据,包括网络设备、网关IP及目的网段
	//相当于 route add -net 0.0.0.0/0 gw {Bridge网桥地址｝ dev ｛容器内的Vet 端点设备｝
	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        endpoint.Network.IpNet.IP,
		Dst:       cidr,
	}
	//调用netlink的RouteAdd,添加路由到容器的网络空间
	//RouteAdd 函数相当于route add命令
	netlink.RouteAdd(defaultRoute)
	return nil
}

func Connect(networkName string, containerInfo *containerInfo.ContainerInfo) error {
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("没有找到匹配的Network: %s", networkName)
	}
	//为容器分配IP地址
	ip, err := ipAddressManage.Allocate(network.IpNet)
	if err != nil {
		return err
	}
	//创建网络端点
	endpoint := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", containerInfo.Id, networkName),
		IPAddress:   ip,
		Network:     network,
		PortMapping: containerInfo.PortMapping,
	}
	//调用网络驱动挂载和配置网络端点
	drivers[network.Driver].Connect(network, endpoint)
	//到容器的namespace配置容器网络设备的IP地址
	configEndpointIpAddressAndRoute(endpoint, containerInfo)
	//配置端口映射信息
	return configPortMapping(endpoint)
}

//配置端口映射
func configPortMapping(endpoint *Endpoint) error {
	//遍历容器端口映射列表
	for _, pm := range endpoint.PortMapping {
		//分成宿主机的端口和容器的端口
		portMapping := strings.Split(pm, ":")
		if len(portMapping) != 2 {
			logrus.Errorf("映射端口格式错误: %v", pm)
			continue
		}
		//由于 iptables没有Go语言版本的实现,所以采用exec.Command 的方式直接调用命令配置
		//在iptables的PREROUTING中添加DNAT规则,将宿主机的端口请求转发到容器的地址和端口上
		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
			portMapping[0], endpoint.IPAddress.String(), portMapping[1])
		//执行 iptables 命令 ， 添加端口映射转发规则
		cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
		output, err := cmd.Output()
		if err != nil {
			logrus.Errorf("iptables Output, %v", output)
			continue
		}
	}
	return nil
}
//func Disconnect(networkName string, cinfo *containerInfo.ContainerInfo) error {
//	return nil
//}
