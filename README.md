# Etpmls-Micro

Englsih | [简体中文](./README_zh-CN.md)

## In Principle:

Etpmls belongs to an organization, not an individual. The project needs more developers to have a future, and we eagerly hope that you can join us.

1. Don't cheat, don't brush star, and would rather nobody use it than fake
2. No cult of personality, every developer is equal, no matter their level is good or bad
3. Swearing is welcome. If you feel that we are not writing well, you can scold it. We like negative comments because it allows us to see ourselves clearly.
4. Developers are highly democratic, voting to determine the future of the project, and the minority obeys the majority, even if you want to `rm -rf /`
5. Born for interest
6. Take open source, repay open source

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

## Configuration

EM needs to configure two files, one is environment variable configuration, the other is application configuration

### Environment Variable Configuration

> .env

```
DEBUG="FALSE"
INIT_DATABASE="FALSE"
```

`DEBUG:`

Whether to enable debugging mode. (TRUE/FALSE), If you fill in TRUE, the **storage/config/app.yaml**  file is read by default, If you fill in FALSE, the **storage/config/app_debug.yaml**  file is read by default

`INIT_DATABASE:`

Whether to initialize the database (TRUE/FALSE),

it is recommended to use it when deploying EA for the first time.

If this mode is turned on, initialization data will be automatically inserted into the database.

Do not turn on this mode when data already exists!

### Application configuration

You need to create two files **app.yaml** (production environment configuration) and **app_debug.yaml** (debug environment configuration) under the storage/config folder. Which file the application uses depends on the value of your environment variable `DEBUG`.

You can refer to the app.yaml.example file to configure.

> The configuration example of Etpmls-Micro/file/app.yaml.example in the EM framework source code is always the latest. If your application configuration file is old, please copy the configuration file example from EM to your project.