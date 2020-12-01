# Etpmls-Micro

[English](./README.md) | 简体中文

## 前提条件
使用前，请确保你是否符合框架所需技能要求
1.具有Protobuf的基础
2.具有Go的基础

## 介绍
Etpmls-Micro（简称EM）是一个微服务框架，使用本框架可以在短时间内快速开发出你的微服务应用。本项目基于Grpc+Grpc Gateway开发。

## 安装
使用go mod安装
```go
import "github.com/Etpmls/Etpmls-Micro"
```

## 快速入门
```go
package main

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

