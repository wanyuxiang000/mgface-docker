package containerNet

import (
	"fmt"
	"net"
	"testing"
)

func TestIPAM_Allocate(t *testing.T){
	_, subnet, _ := net.ParseCIDR("192.168.254.0/23")
	ip,_:= ipAddressManage.Allocate(subnet)
	fmt.Println(ip)
}

func TestIPAM_Release(t *testing.T) {
	_, subnet, _ := net.ParseCIDR("192.168.254.0/23")
	ip:=net.ParseIP("192.168.254.1")
	ipAddressManage.Release(subnet,&ip)
}