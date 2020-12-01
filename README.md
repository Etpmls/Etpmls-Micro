# Etpmls-Micro

Englsih | [简体中文](./README_zh-CN.md)

## Prerequisites
Before using, please make sure you meet the required skills of the framework
1. Have the foundation of Protobuf
2. Have a foundation of Go

## Introduction
Etpmls-Micro (EM for short) is a micro-service framework, using this framework can quickly develop your micro-service applications in a short time.This project is developed based on Grpc+Grpc Gateway.

## Installation
Install with go mod
```go
import "github.com/Etpmls/Etpmls-Micro"
```

## Quick start
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
This is the simplest application. You only need to implement three methods to successfully register a microservice application.

`RegisterRpcService`: Implement RPC service

`RegisterHttpService`: Implement HTTP service

`RegisterRoute`: Implement routing