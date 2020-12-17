# Etpmls-Micro

Englsih | [简体中文](./README_zh-CN.md)

## In Principle:

Etpmls belongs to an organization, not an individual. The project needs more developers to have a future, and we eagerly hope that you can join us.

1. Don't cheat, don't brush star, and would rather nobody use it than fake
2. Swearing is welcome. If you feel that we are not writing well, you can scold it. We like negative comments because it allows us to see ourselves clearly.
3. Developers are highly democratic, voting to determine the future of the project, and the minority obeys the majority, even if you want to `rm -rf /`
4. Born for interest

## Prerequisites

Before using, please make sure you meet the required skills of the framework

1. Have the foundation of Protobuf

2. Have a foundation of Go

## Introduction
Etpmls-Micro (EM for short) is a micro-service framework, using this framework can quickly develop your micro-service applications in a short time.This project is developed based on Grpc+Grpc Gateway.

>We recommend that you use the following items together to facilitate rapid development of your application.
>
>[EM-Auth](https://github.com/Etpmls/EM-Auth): The main control center, which integrates RBAC0 authentication of users, roles and permissions, custom menus, cache clearing, disk cleaning and other functions
>
>[EM-Attachment](https://github.com/Etpmls/EM-Attachment): Attachment center, used to process attachments of various microservices.

## Installation
Install with go mod
```go
import "github.com/Etpmls/Etpmls-Micro"
```

## Quick start
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
This is the simplest application. You only need to implement three methods to successfully register a microservice application.

`RegisterRpcService`: Implement RPC service

`RegisterHttpService`: Implement HTTP service

`RegisterRoute`: Implement routing