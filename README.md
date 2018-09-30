# mgface-docker

#### 项目介绍
mgface-docker golang实现docker的功能

#### 软件架构
├── aufs
│   ├── deleteFileSystem.go
│   └── newFileSystem.go
├── cgroup
│   ├── cgroupManager.go
│   ├── cgroupTools.go
│   ├── cpu.go
│   ├── cpuset.go
│   ├── memory.go
│   └── resouceConfig.go
├── command
│   ├── commitCommand.go
│   ├── execCommand.go
│   ├── initCommand.go
│   ├── listCommand.go
│   ├── logCommand.go
│   ├── networkCommand.go
│   ├── rmCommand.go
│   ├── runCommand.go
│   ├── startCommand.go
│   └── stopCommand.go
├── constVar
│   └── constVariables.go
├── container
│   ├── commit.go
│   ├── exec.go
│   ├── init.go
│   ├── logs.go
│   ├── ps.go
│   ├── rm.go
│   ├── run.go
│   ├── start.go
│   └── stop.go
├── containerInfo
│   ├── containerInfo.go
│   ├── randStringBuffer.go
│   └── randStringBuffer_test.go
├── containerNet
│   ├── bridgeDriver.go
│   ├── driver.go
│   ├── golang daemon实现.md
│   ├── hostPortUp.go
│   ├── init.go
│   ├── ipam.go
│   ├── ipam_test.go
│   ├── network.go
│   └── networkTools.go
├── Gopkg.lock
├── Gopkg.toml
├── main.go
├── nsenter
│   └── setns.go


#### 安装教程

1. 首先要下载golang并且安装
2. 下载该代码
3. 设置环境变量
   export GOROOT=/usr/local/go
   export GOPATH=/usr/local/goproject/mgface-docker
   export APP=$GOPATH/bin
   export PATH=$APP:$PATH:$GOPATH:$GOROOT/bin

#### 使用说明

1. mgface.com --help

#### 参与贡献

1. Fork 本项目
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request


#### 码云特技

1. 使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2. 码云官方博客 [blog.gitee.com](https://blog.gitee.com)
3. 你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解码云上的优秀开源项目
4. [GVP](https://gitee.com/gvp) 全称是码云最有价值开源项目，是码云综合评定出的优秀开源项目
5. 码云官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6. 码云封面人物是一档用来展示码云会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)