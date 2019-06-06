# Golang微服务开发实践

微服务概念学习：可参考 [Nginx 的微服务文章](https://www.nginx.com/blog/introduction-to-microservices/)

微服务最佳实践：可参考 [微服务最佳实践](https://juejin.im/post/5cbbe051f265da03973aabcb#heading-24)

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

首先初始化：`make init`  

**Makefile** 部分代码如下：

```makefile
init:
    cd ..
	mv MicroServicePractice ${GOPATH}/src/Ethan/
	./pull.sh # 安装 go 依赖
	cd plugins
	docker-compose -f docker-compose.yml up -d # 安装插件，如：kafka, consul, zookeeper, jaeger
```

之后就可以运行代码了：

注：建议自己开多个终端 `go run` ，这样可以看日志

```shell
make run # 允许 服务 server
```

测试：

注：注意顺序，刚开始啥数据都没有的

```shell
go run user-cli/cli.go
export Token=$Token # 注意换成前面生成的Token
go run vessel-cli/cli.go
go run consignment-cli/cli.go
```

![peek](dist/peek.gif)



![consul](/home/mu-mo/桌面/Go/src/github.com/yun-mu/MicroServicePractice/dist/consul.png)



![jaeger](/home/mu-mo/桌面/Go/src/github.com/yun-mu/MicroServicePractice/dist/jaeger.png)



![mongo2](/home/mu-mo/桌面/Go/src/github.com/yun-mu/MicroServicePractice/dist/mongo2.png)

![mongo1](/home/mu-mo/桌面/Go/src/github.com/yun-mu/MicroServicePractice/dist/mongo1.png)

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

这里使用 micro 插件，若想和不使用插件对比，可使用如下命令：

```shell
protoc --proto_path=proto:. --go_out=out/ --micro_out=out/ proto/vessel/vessel.proto
```

这样会生成两个文件，一个为 *.micro.go* 一个为 *.pb.go*

这里顺便看一下 生成的 pb 文件里是如何进行 rpc 调用的，我们随便看一个 方法，如：*vessel* 的 *FindAvailable*

```go
func (c *vesselServiceClient) FindAvailable(ctx context.Context, in *Specification, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "VesselService.FindAvailable", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
```



### 微服务开发流程

如果使用 *grpc* 作为 *server* 和 *client*，开发流程如下：

注：*server* 和 *client* 必须相同，如：我的代码中 *server* 和 *client* 使用的都是 *rpc*, *transport* 是 *tcp* 

![image-20180512044329199](https://images.yinzige.com/2018-05-11-204329.png)

### 目录简介

- api：对外暴露的HTTP web 接口，可以理解为 网关
- common：所有服务都能调用的东西，如 *GetMicroClient, GetMicroServer* 
- config：配置中心，其他服务的启动都依赖的配置
- consignment
- consignment-cli：cli 测试
- interface-center：proto 文件中心，同时生成的 .go 文件也在这里
- shippy-ui：前端测试 ui 代码，对接API，API还没写完
- user
- user-cli
- vessel
- vessel-cli



## go-micro 详解

*micro* 文档：https://micro.mu/docs/index.html

参见另一篇 [go-micro详解](https://github.com/yun-mu/MicroServicePractice/blob/master/go-micro%E8%AF%A6%E8%A7%A3.md)

