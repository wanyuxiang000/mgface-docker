package containerNet

import (
	"fmt"
	"net"
	"testing"
)

func TestAllocate(t *testing.T){
	_, subnet, _ := net.ParseCIDR("172.168.0.0/16")
	ip,_:=ipAllocator.Allocate(subnet)
	fmt.Println(ip)
}