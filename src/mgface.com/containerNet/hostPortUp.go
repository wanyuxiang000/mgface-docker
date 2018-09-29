package containerNet

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"net"
)

/*
#include <unistd.h>
*/
import "C"

//houstport是宿主机的端口,作为对外代理的端口
//containerIp 容器的IP地址和端口
func hostServer(hostport string, containerIp string,tty bool) {
	//假如没开启tty终端的话，那么是后台启动
	//if !tty{
	//	logrus.Info("通过Cgo来实现监听tcp的daemon实现")
	//	//设置守护进程
	//	C.daemon(1, 1)
	//	runtime.GOMAXPROCS(runtime.NumCPU())
	//}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", hostport))
	if err != nil {
		logrus.Infof("监听宿主机端口报错:%s,错误信息:%+v", hostport, err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Infof("建立连接错误:%+v", err)
			continue
		}
		logrus.Infof("远程地址:%+s,本地地址:%s", conn.RemoteAddr(), conn.LocalAddr())
		go handle(conn, containerIp)
	}
}

func handle(sconn net.Conn, ip string) {
	defer sconn.Close()
	dconn, err := net.Dial("tcp", ip)
	if err != nil {
		logrus.Errorf("连接%+v失败:%+v\n", ip, err)
		return
	}
	ExitChan := make(chan bool, 1)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		_, err := io.Copy(dconn, sconn)
		logrus.Infof("往%+v发送数据失败:%+v\n", ip, err)
		ExitChan <- true
	}(sconn, dconn, ExitChan)


	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		_, err := io.Copy(sconn, dconn)
		logrus.Infof("从%+v接收数据失败:%+v\n", ip, err)
		ExitChan <- true
	}(sconn, dconn, ExitChan)

	<-ExitChan
	dconn.Close()
}