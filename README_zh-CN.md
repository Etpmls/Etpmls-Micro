# Etpmls-Micro

[English](./README.md) | 简体中文

##  原则

Etpmls属于一个组织，而不是个人，项目需要更多的开发者才能有未来，迫切的希望您能够加入我们。

1. 不作弊，不刷star，宁可没人使用，也不去造假
2. 不搞个人崇拜，每个开发者都是平等的，无论他们的水平是好是坏
3. 欢迎谩骂。如果你觉得我们哪里写的不好，可以骂出来，我们喜欢负面的评价，因为这样可以使我们看清自己。
4. 开发者高度民主，投票决定项目的未来，少数服从多数，哪怕你想`rm -rf /`
5. 为兴趣而生
6. 取用开源，回报开源

## 前提条件

使用前，请确保你是否符合框架所需技能要求

1. 具有Protobuf的基础
2. 具有Go的基础
3. 具有Docker的基础

## 介绍

Etpmls-Micro（简称EM）是一个微服务框架，使用本框架可以在短时间内快速开发出你的微服务应用。本项目基于Grpc+Grpc Gateway开发。

>我们推荐您搭配以下项目同时使用，便于快速开发您的应用。
>
>[EM-Auth](https://github.com/Etpmls/EM-Auth)：总控制中心，集成用户、角色、权限的RBAC0的鉴权、自定义菜单、清除缓存、磁盘清理等功能
>
>[EM-Attachment](https://github.com/Etpmls/EM-Attachment)： 附件中心，用于处理各个微服务的附件。

![Process](docs/images/Process.jpg)

## 版本说明

我们的版本格式为：vA.B.C，

如果你的EM版本与最新版本只有`C`不同，那么你无需犹豫，可以直接升级。

如果你的EM版本与最新版本只有`B`不同，那么你可能需要关注一下升级手册，因为我们可能有一些改动，当然我们也会努力控制高版本的兼容性。

如果你的EM版本与最新版本只有`A`不同，那么说明版本更新较大或重构，你要考虑是否要升级到最新的版本。

> 我们尽量以高兼容性、小改动为主，降低使用者的学习成本

## 安装
使用go mod安装
```go
import "github.com/Etpmls/Etpmls-Micro"
```

## 快速入门

### 说明

```go
// Etpmls-Micro/example/main.go
package main

import "github.com/Etpmls/Etpmls-Micro"

func main()  {
	var reg = em.Register{
		GrpcServiceFunc: RegisterRpcService,
		HttpServiceFunc: RegisterHttpService,
		RouteFunc:       RegisterRoute,
	}
	reg.Init()
	reg.Run()
}
```
这个是最简单的应用，你只需要实现三个方法，即可成功注册一个微服务应用。

`RegisterRpcService` ： 实现RPC服务

`RegisterHttpService` ： 实现HTTP服务

`RegisterRoute` : 实现路由

### 运行

EA目录中已经包含了示例，这是一个最精简的示例，以便于你的理解。

进入Etpmls-Micro/example并且执行

```shell
go run main.go
```

在浏览器输入http://localhost:8081/hello ，你会发现示例返回world。

> 这个示例没有包含任何服务，仅仅只有一个HTTP路由。

## 目录规范

> 请注意，storage目录下要完全符合EM的命名规范，对于其他目录的命名要求并不是强制性的，如果您不喜欢我们的目录规范，可以跳过本章。

我们并没有严格的目录规范，您可以完全定义您的目录名（storage目录除外），但是我们推荐您使用我们的目录规范，目录名的标准化有助于提高阅读体验，提升其他开发者的阅读效率。

```
# storage目录要符合EM的命名规范
/storage
|______/config
|______/language
|______/log
|______/menu
|______/upload

# 下面的目录命名我们并不强制，只是推荐您这样使用
/src
|______/application
|______|______/client
|______|______/database
|______|______/middleware
|______|______/model
|______|______/protobuf
|______|______|______proto
|______|______/service

|______/register
```

`/storage` : **[名称不可修改]** 存放无需编译的文件的目录

`/storage/config` : **[名称不可修改]** 存放应用配置文件

`/storage/language` : **[名称不可修改]** 存放多语言文件

`/storage/log` : **[名称不可修改]** 存放日志文件

`/storage/menu` : **[名称不可修改]** 存放前端自定义菜单的json文件，通常EM-Auth需要此目录，其他微服务应用可不创建此目录

`/storage/upload` : 存放上传附件的目录，通常EM-Attachment需要此目录，其他微服务应用可不创建此目录

> 下面的目录命名我们并不强制，只是推荐您这样使用

`/src` : 源代码放在这个目录

`/src/application` : 业务逻辑、功能

`/src/application/client` : 与其他微服务交互，作为客户端请求其他微服务

`/src/application/database` : 数据库字段定义

`/src/application/middleware` : 自定义中间件

`/src/application/model` : 模型，处理服务的业务逻辑

`/src/application/protobuf` : Protobuf编译后的文件存放位置

`/src/application/protobuf/proto` : Protobuf编译前的文件存放位置

`/src/application/service`  : 相当于MVC中的controller，处理服务请求

`/src/application/service.go`  : 定义该服务特有的一些声明等（如常量、变量）

`/src/register` :注册逻辑（如注册路由、注册中间件、注册数据库等）

## 环境搭建

本项目需要结合Traefik（网关）和Consul（服务发现），仅贴出搭建相关文件代码，并不会详细说明细节原理，具体的细节可以参考官方文档。[Traefik](https://doc.traefik.io/)  [Consul](https://www.consul.io/docs)

创建一个文件夹存放docker-compose，本文以`.`代表当前目录

> ./docker-compose.yml

```yaml
version: '3'
    
services:
  traefik:
    # The official v2 Traefik docker image
    image: traefik:v2.4
    environment:
    - TZ=Asia/Shanghai
    # Enables the web UI and tells Traefik to listen to docker
    command: --api.insecure=true --providers.docker
    ports:
      # The HTTP port
      - "80:80"
      # The Web UI (enabled by --api.insecure=true)
      - "8080:8080"
      - "443:443"
    restart: on-failure
    volumes:
      # So that Traefik can listen to the Docker events
      - /var/run/docker.sock:/var/run/docker.sock
      - ./traefik:/etc/traefik
    networks:
      em:
       ipv4_address: [##定义traefik的IP地址，如10.0.0.2##]
      
  consul:
    image: consul:1.8.5
    volumes:
      - ./consul/data:/consul/data
    ports:
      - "8300:8300/tcp"
      - "8301:8301/udp"
      - "8302:8302/udp"
      - "8500:8500/tcp"
      - "53:8600/udp"
    restart: on-failure
    networks:
      em:
       ipv4_address: [##定义consul的IP地址，如10.0.0.2##]
       
networks:
 em:
  ipam:
   driver: default
   config:
    - subnet: "[##定义网关地址，如10.0.0.0/24##]"
```

> ./trafik/traefik.yaml

```yaml
api:
  dashboard: true
 
entryPoints:
  web:
    address: ":80"
    http:
      redirections:
        entryPoint:
          to: websecure
          scheme: https
  websecure:
    address: ":443"
    forwardedHeaders:
     trustedIPs:
     - "[##填写你定义的网关地址，如10.0.0.0/24##]"
  dashboard:
    address: ":8080"

providers:
  file:
    directory: /etc/traefik/config
  consulCatalog:
   refreshInterval: 30s
   prefix: em
   endpoint:
    address: [##定义consul的地址和端口，如10.0.0.3:8500##]


certificatesResolvers:
  myresolver:
    acme:
      email: [##邮箱，获取SSL证书用到##]
      storage: acme.json
      httpChallenge:
        entryPoint: web
        
log:
  filePath: "/etc/traefik/log/error.log"
  format: "json"
  level: WARN
accessLog:
  filePath: "/etc/traefik/log/access.log"
  format: "json"
  bufferingSize: 100
```

> ./traefik/config/dashboard.yaml

```yaml
http:
 routers:
  dashboard:
   entryPoints:
   - "dashboard"
   rule: (PathPrefix(`/api`) || PathPrefix(`/dashboard`))
   service: api@internal
   middlewares:
    - auth

 middlewares:
  auth:
   basicAuth:
    users:
    - "admin:$apr1$sadEhKwW$BNpyOakcbLp/P7JyP5ghs0"     # admin admin
```

>  ./traefik/config/em.yaml

```yaml
http:
 routers:
  em-template:
    entryPoints:
    - "web"
    - "websecure"
    rule: "Host(`[##填写网站域名，如baidu.com##]`)"
    service: em-template
    middlewares:
     - rateLimit
    tls: 
     certResolver: myresolver

 middlewares:
  forwardAuth:
   forwardAuth:
    address: "[##填写traefik定义的地址+权限验证的地址，如https://10.0.0.2/api/checkAuth##]"
    tls:
     insecureSkipVerify: true
  rateLimit:
   rateLimit:
    average: 1000
    period: 10s
    burst: 2000
  circuitBreaker_em-auth:
   circuitBreaker:
    expression: "NetworkErrorRatio() > 0.30 || LatencyAtQuantileMS(50.0) > 3000"
  circuitBreaker_em-attachment:
   circuitBreaker:
    expression: "NetworkErrorRatio() > 0.30 || LatencyAtQuantileMS(50.0) > 3000"
       
 services:
  em-template:
   loadBalancer:
    servers:
    - url: "[##定义前端地址，如http://192.168.3.225:9527/##]"
```

运行docker-compose

```shell
docker-compose up -d
```

## 配置

### EM配置

EM需要配置两个文件，一个是环境变量配置，另一个是应用配置

#### 环境变量配置

> .env

```
DEBUG="FALSE"
INIT_DATABASE="FALSE"
```

你可以参考.env.example文件来配置。

`DEBUG:`

是否开启调试模式。(TRUE/FALSE)，

若填写FALSE，则默认读取配置文件**storage/config/app.yaml**， 若填写TRUE，则默认读取配置文件**storage/config/app_debug.yaml**

`INIT_DATABASE:`

> 默认为False，如果你想开启此功能，你需要在em.Register中注册DatabaseMigrate

是否初始化数据库(TRUE/FALSE)，

建议第一次部署EA时使用。

如开启此模式，将自动向数据库中插入初始化数据。

请勿在已存在数据的情况下开启此模式！

#### 应用配置

你需要在storage/config文件夹下方建立两个文件**app.yaml**（生产环境配置）和**app_debug.yaml**（debug环境配置），应用是否使用哪个文件取决于你的环境变量`DEBUG`的值。

你可以参考app.yaml.example文件来配置。

> EM框架源码中的Etpmls-Micro/file/app.yaml.example配置示例永远是最新的。如果您打算从低版本升级更高版本的EM框架，请从EM中复制最新的配置文件示例到你的项目下。

### 网关配置

本项目网关以Traefik为例、服务发现以Consul为例。若想整合网关与服务发现，需要进行配置网关。

> Traefik官方参考文章
>
> https://doc.traefik.io/traefik/providers/consul-catalog/
>
> https://doc.traefik.io/traefik/routing/providers/consul-catalog/

我们需要把相关配置写在`storage/config/app[_debug].yaml`文件的service-discovery.service.rpc/http.tag中，我们提供一个示例参考，您可以在此基础上直接修改。

```yaml
      tag: [
        "em.http.routers.[YOUR_SERVICE_NAME].entrypoints=web,websecure",
        "em.http.routers.[YOUR_SERVICE_NAME].rule=Host(`[YOUR_DOMAIN]`) && PathPrefix(`[YOUR_SERVICE_ROUTE_PATH]`)",
        "em.http.routers.[YOUR_SERVICE_NAME].tls.certresolver=myresolver",
        "em.http.routers.[YOUR_SERVICE_NAME].middlewares=circuitBreaker_em-attachment@file,forwardAuth@file",
        "em.http.routers.[YOUR_SERVICE_NAME].service=[YOUR_SERVICE_NAME]",

        "em.http.services.[YOUR_SERVICE_NAME].loadbalancer.passhostheader=true",
      ]
```

> [YOUR_DOMAIN]

替换为你的域名

> [YOUR_SERVICE_ROUTE]

替换为你的服务路径，如/api/attachment/

> [YOUR_SERVICE_NAME]

替换为你的服务名