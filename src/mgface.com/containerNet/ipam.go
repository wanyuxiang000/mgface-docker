package containerNet

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	. "mgface.com/constVar"
	"net"
	"os"
	"path"
	"strings"
)

//存放IP地址分配信息
type IPAM struct {
	//分配文件存放位置
	SubnetAllocatorPath string
	//key是网段,value是分配的位图数组
	Subnets *map[string]string
}

var ipAddressManage = &IPAM{
	SubnetAllocatorPath: IpamDefaultAllocatorPath,
}

func (ipam *IPAM) load() error {
	logrus.Infof("加载ipam文件:%s",ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	subnetConfigFile, _ := os.Open(ipam.SubnetAllocatorPath)
	defer subnetConfigFile.Close()
	subnetJson := make([]byte, 1024*1024)
	n, _ := subnetConfigFile.Read(subnetJson)
	json.Unmarshal(subnetJson[:n], ipam.Subnets)
	return nil
}

func (ipam *IPAM) dump() error {
	ipamConfigDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipamConfigDir, 0644)
		} else {
			return err
		}
	}
	//打开存储文件,os.O_TRUNC表示如果存在则清空os.O_CREATE表示如果不存在则创建
	subnetConfigFile, _ := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	defer subnetConfigFile.Close()
	content, _ := json.Marshal(ipam.Subnets)
	content = append(content,[]byte("\n")...)
	subnetConfigFile.Write(content)
	return nil
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	//存放网段中地址分配信息的数组
	ipam.Subnets = &map[string]string{}
	// 从文件中加载已经分配的网段信息
	ipam.load()
	_, subnet, _ = net.ParseCIDR(subnet.String())
	//net.IPNet.Mask.Size()函数会返回网段的子网掩码的总长度和网段前面的固定位的长度
	//比如“127.0.0.1/8”网段的子网掩码是"255.0.0.0"
	//那么Mask.Size()的返回值就是前面 255 所对应的位数和总位数,即8和32
	ones, bits := subnet.Mask.Size()

	ipBits, exist := (*ipam.Subnets)[subnet.String()]
	//如果之前没有分配过这个网段,则初始化网段的分配配置
	if !exist {
		//用"0"填满这个网段的配置,1<<uint8(bits-ones)表示这个网段中有多少个可用地址
		//bits-ones是子网掩码后面的网络位数， 2^(bits-ones)表示网段中的可用IP数
		//而 2^(bits-ones)等价于1<<uint8(bits-ones)
		ipBits = strings.Repeat("0", 1<<uint8(bits-ones))
	}

	for c := range ipBits {
		if ipBits[c] == '0' {
			//Go的字符串,创建之后就不能修改,所以通过转换成byte数组,修改后再转换成字符串赋值
			ipalloc := []byte(ipBits)
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)
			//这里的IP为初始IP,比如对于网段192.168.0.0/16,这里就是192.168.0.0
			ip = subnet.IP
			//通过网段的IP与上面的偏移相加计算出分配的IP地址，由于IP地址是uint的一个数组，
			//需要通过数组中的每一项加所需要的值，比如网段是172.16.0.0/12,数组序号是65555.
			//那么在［172.16.0.0］上依次加［uint8(65555>>24)、 uint8(65555>>16)、
			//uint8(65555>>8)、 uint8(65555>>0),即［0,1,0,19]那么获得的IP就是172.17.0.19.
			for t := uint(4); t > 0; t -= 1 {
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			//由于此处IP是从1开始分配的，所以最后再加1，最终得到分配的IP是172.17.0.20
			ip[3] += 1
			break
		}
	}
	ipam.dump()
	return
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}
	_, subnet, _ = net.ParseCIDR(subnet.String())
	ipam.load()
	c := 0
	//将IP地址转换成4个字节的表示方式
	releaseIP := ipaddr.To4()
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t -= 1 {
		/*与分配IP相反,释放IP获得索引的方式是IP地址的每一位相减之后分别左移将对应的数值加
		到索引上。*/
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}
	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)
	ipam.dump()
	return nil
}
