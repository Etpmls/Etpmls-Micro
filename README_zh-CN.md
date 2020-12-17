# Etpmls-Micro

[English](./README.md) | 简体中文

## 原则：
Etpmls属于一个组织，而不是个人，项目需要更多的开发者才能有未来，迫切的希望您能够加入我们。

1. 不作弊，不刷star，宁可没人使用，也不去造假
2. 不搞个人崇拜，每个开发者都是平等的，无论他们的水平是好是坏
3. 欢迎谩骂。如果你觉得我们哪里写的不好，可以骂出来，我们喜欢负面的评价，因为这样可以使我们看清自己。
4. 开发者高度民主，投票决定项目的未来，少数服从多数，哪怕你想`rm -rf /`
5. 为兴趣而生
6. 取用开源，回报开源

## 前提条件

使用前，请确保你是否符合框架所需技能要求

1.具有Protobuf的基础

2.具有Go的基础

## 介绍
Etpmls-Micro（简称EM）是一个微服务框架，使用本框架可以在短时间内快速开发出你的微服务应用。本项目基于Grpc+Grpc Gateway开发。

>我们推荐您搭配以下项目同时使用，便于快速开发您的应用。
>
>[EM-Auth](https://github.com/Etpmls/EM-Auth)：总控制中心，集成用户、角色、权限的RBAC0的鉴权、自定义菜单、清除缓存、磁盘清理等功能
>
>[EM-Attachment](https://github.com/Etpmls/EM-Attachment)： 附件中心，用于处理各个微服务的附件。

## 安装
使用go mod安装
```go
import "github.com/Etpmls/Etpmls-Micro"
```

## 快速入门
```go
package main

import "github.com/Etpmls/Etpmls-Micro"

func main() {
	var reg = em.Register{
		GrpcServiceFunc:    	RegisterRpcService,
		HttpServiceFunc:    	RegisterHttpService,
		RouteFunc:          	RegisterRoute,
	}

	reg.Run()
}
```
这个是最简单的应用，你只需要实现三个方法，即可成功注册一个微服务应用。

`RegisterRpcService` ： 实现RPC服务

`RegisterHttpService` ： 实现HTTP服务

`RegisterRoute` : 实现路由

