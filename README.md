# Microservice

- [Microservice](#microservice)
  - [概述](#概述)
    - [依赖](#依赖)
  - [安装](#安装)
    - [安装 Mysql](#安装-mysql)
    - [安装 Redis](#安装-redis)
    - [安装 Consul](#安装-consul)
    - [安装 Microservice](#安装-microservice)
  - [配置](#配置)
  - [协议结构](#协议结构)
    - [to server msg](#to-server-msg)
    - [to client msg](#to-client-msg)
  


## 概述
  这是一个基于golang开发的游戏微服务架构，每个服务都实现了水平扩展。
![RUNOOB 图标](https://github.com/stringMao/SrcRepository/blob/main/%E5%BE%AE%E6%9C%8D%E5%8A%A1%E6%9E%B6%E6%9E%84.png)

### 依赖
  本架构使用了一些第三方开源库来实现相关功能
- HTTP接口(github.com/gin-gonic/gin)        
- 数据库连接 (github.com/go-xorm/xorm)
- Redis连接 (github.com/gomodule/redigo/redis)
- 协议(github.com/golang/protobuf)
- 日志(github.com/sirupsen/logrus)
- 服务注册(github.com/hashicorp/consul)
- 配置读取(github.com/Unknwon/goconfig)

## 安装

### 安装 Mysql
 - [mysql安装教程](https://www.runoob.com/mysql/mysql-install.html)

 - 导入脚本
### 安装 Redis
 - [redis安装教程](https://www.runoob.com/redis/redis-install.html)
### 安装 Consul
- [consul下载](https://www.consul.io/downloads.html )
### 安装 Microservice

```
git clone https://github.com/stringMao/Microservice.git

cd LoginSvr

go build

cd ../GateSvr

go build

cd ../HallSvr

go build

```

## 配置
每个服务器启动前，同级目录下都需要配置app.ini。例如以下是GateSvr的配置模板

```
[log]
level = debug
[webmanager]
ip=127.0.0.1
port =8092
# 服务发现consul的地址
consuladdr = 127.0.0.1:8500


[server]
#该类型服务下的唯一ID
sid = 1


[business]
#客户端连接监听端口
clientPort=8090
#服务器连接的监听端口
serverPort=8091


[redis]
host = 127.0.0.1
port = 6379
username = root
password = ***
database=1
maxopenconns = 20
maxidleconns=5

[mysql]
host = 127.0.0.1
port = 3306
username = ###
password = ***
dbname = player
maxopenconns = 100
maxidleconns=10
```

## 协议结构
###  msg struct
![RUNOOB 图标](https://github.com/stringMao/SrcRepository/blob/main/%E5%BE%AE%E6%9C%8D%E5%8A%A1/msgStruct.png)



