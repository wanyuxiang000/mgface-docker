package containerNet

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
	"os"
	"path"
)

var (
	drivers  = map[string]Driver{}
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
	content := make([]byte, 1024*1024)
	n, _ := configFile.Read(content)
	if err := json.Unmarshal(content[:n], network); err != nil {
		logrus.Infof("解析文件 [%s] 失败,请核实信息:%+v.", dumpPath, err)
	}
	return nil
}
