https://www.jianshu.com/p/1d680721780f

实现tcp server的时候发现官方没有提供daemon的方式，在网上搜了一下，大概有下面几个方案：

1.nohup 2.supervise 3.Cgo deamon函数 4.go通过syscall调用fork实现(这个和第3条原理一样)

code 每种方式都各有优劣，这里说一下通过Cgo如何实现daemon：

package main

import (
"fmt"
"net"
"runtime"
)

/*

# include <unistd.h>

*/ import "C"

func main() { // 守护进程 C.daemon(1, 1)
runtime.GOMAXPROCS(runtime.NumCPU())
fmt.Println("Starting the server ...")
listener, err := net.Listen("tcp", "localhost:8080")
if err != nil { fmt.Println("Error listening", err.Error())
return }

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("error accepting", err.Error())
            return
        }

        go doServerStuff(conn)
    }

}

func doServerStuff(conn net.Conn) { for { buf := make([]byte, 512)
_, err := conn.Read(buf)
if err != nil { fmt.Println("Error reading", err.Error())
return } fmt.Printf("Received data: %v", string(buf))
} } 验证 go install daemon ./bin/daemon

telnet 127.0.0.1 8080
