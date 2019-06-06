# Golang微服务开发实践

微服务概念学习：可参考 [Nginx 的微服务文章](https://www.nginx.com/blog/introduction-to-microservices/)

## demo 简介

服务：

- consignment-service（货运服务）
- user-service（用户服务）
- log-service (日志服务)
- vessel-service（货船服务）
- api-service (API 服务)

用到的技术栈如下：

```yaml
framework: go-micro, gin
Transport: tcp
Server: rpc
Client: rpc
RegisterTTL: 30s
RegisterInterval: 20s
Registry: consul, 服务发现和注册
Broker: kafka, 消息队列
Selector: cache, 负载均衡
Codec: protobuf, 编码
Tracing: jaeger, 链路追踪
Metrics: jaeger
breaker: hystrix, 熔断
ratelimit: uber/ratelimit, 限流
```

### 服务关系图

![project](dist/project.png)

### 实体关系图

![image-20180512010554833](https://images.yinzige.com/2018-05-11-170555.png)

### 服务流程示例

![image-20180522174448548](https://images.yinzige.com/2018-05-22-094448.png)

### 认证

采用 JWT

![image-20180528105549866](https://images.yinzige.com/2018-05-28-025550.png)

### 发布订阅模式

![未命名文件](dist/kafka.png)

### demo 运行

前提工具：`go, dep, docker, docker-compose, mongo`

首先初始化：`make init`  注：建议阅读**Makefile**文件

之后就可以运行代码了：

注：建议自己开多个终端 `go run` ，这样可以看日志

```shell
make run
```

测试：

注：注意顺序，刚开始啥数据都没有的

```shell
go run user-cli/cli.go
export Token=$Token # 注意换成前面生成的Token
go run vessel-cli/cli.go
go run consignment-cli/cli.go
```

## 开发详解

### proto 代码生成

安装工具：

*protoc* 安装：http://google.github.io/proto-lens/installing-protoc.html

*protoc-gen-go* 和 *protoc-gen-micro*

```shell
go get -u -v google.golang.org/grpc				
go get -u -v github.com/golang/protobuf/protoc-gen-go
go get -u -v github.com/micro/protoc-gen-micro
```

生成的脚本我已经写好 *Makefile*, 进入 *interface-center* 目录，执行`make build` 即可

内部示例如下：

```shell
protoc --proto_path=proto:. --go_out=plugins=micro:out/ proto/vessel/vessel.proto
```

### 微服务开发流程

如果使用 *grpc* 作为 *server* 和 *client*，开发流程如下：

注：*server* 和 *client* 必须相同，如：我的代码中 *server* 和 *client* 使用的都是 *rpc*, *transport* 是 *tcp* 

![image-20180512044329199](https://images.yinzige.com/2018-05-11-204329.png)



## go-micro 详解

*micro* 文档：https://micro.mu/docs/index.html

参见另一篇 [go-micro详解](https://github.com/yun-mu/MicroServicePractice/blob/master/go-micro%E8%AF%A6%E8%A7%A3.md)

