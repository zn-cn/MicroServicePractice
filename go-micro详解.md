# go-micro 详解

## 整体架构介绍

### 通信流程

​     go-micro的通信流程大至如下

![img](https://user-gold-cdn.xitu.io/2019/2/26/16927b99e35369d8?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

 

​    Server监听客户端的调用，和Brocker推送过来的信息进行处理。并且Server端需要向Register注册自己的存在或消亡，这样Client才能知道自己的状态。

​    Register服务的注册的发现。

​    Client端从Register中得到Server的信息，然后每次调用都根据算法选择一个的Server进行通信，当然通信是要经过编码/解码，选择传输协议等一系列过程的。

​    如果有需要通知所有的Server端可以使用Brocker进行信息的推送。

​    Brocker 信息队列进行信息的接收和发布。

 

​     go-micro之所以可以高度订制和他的框架结构是分不开的，go-micro由8个关键的interface组成，每一个interface都可以根据自己的需求重新实现，这8个主要的inteface也构成了go-micro的框架结构。

![img](https://camo.githubusercontent.com/9057599d2bc2d3c79c43423521d71f4ea0851457/68747470733a2f2f6d6963726f2e6d752f646f63732f696d616765732f676f2d6d6963726f2e737667)

![go-micro](https://github.com/hb-go/micro/raw/master/doc/img/micro.jpg)

这些接口go-micir都有他自己默认的实现方式，还有一个go-plugins是对这些接口实现的可替换项。你也可以根据需求实现自己的插件

![img](https://user-gold-cdn.xitu.io/2019/2/26/16927b99e39b03d1?imageView2/0/w/1280/h/960/format/webp/ignore-error/1) 

### Transport

服务之间通信的接口。也就是服务发送和接收的最终实现方式，是由这些接口定制的。

##### Interface
```go
type Message struct {
	Header map[string]string
	Body   []byte
}

type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
}

type Client interface {
	Socket
}

type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}

// Transport is an interface which is used for communication between
// services. It uses socket send/recv semantics and had various
// implementations {HTTP, RabbitMQ, NATS, ...}
type Transport interface {
	Dial(addr string, opts ...DialOption) (Client, error)
	Listen(addr string, opts ...ListenOption) (Listener, error)
	String() string
}
```

##### Options
```go
type Options struct {
	Addrs     []string
	Codec     codec.Codec
	Secure    bool
	TLSConfig *tls.Config
	// Timeout sets the timeout for Send/Recv
	Timeout time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type DialOptions struct {
	Stream  bool
	Timeout time.Duration

	// TODO: add tls options when dialling
	// Currently set in global options

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type ListenOptions struct {
	// TODO: add tls options when listening
	// Currently set in global options

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
```
​    Transport 的Listen方法是一般是Server端进行调用的，他监听一个端口，等待客户端调用。

​    Transport 的Dial就是客户端进行连接服务的方法。他返回一个Client接口，这个接口返回一个Client接口，这个Client嵌入了Socket接口，这个接口的方法就是具体发送和接收通信的信息。

​    http传输是go-micro默认的同步通信机制。当然还有很多其他的插件：grpc,nats,tcp,udp,rabbitmq,nats，都是目前已经实现了的方式。在go-plugins里你都可以找到。

### Codec

​     有了传输方式，下面要解决的就是传输编码和解码问题，go-micro有很多种编码解码方式，默认的实现方式是protobuf,当然也有其他的实现方式，json、protobuf、jsonrpc、mercury等等。

源码

```
type Codec interface {
    ReadHeader(*Message, MessageType) error
    ReadBody(interface{}) error
    Write(*Message, interface{}) error
    Close() error
    String() string
}

type Message struct {
    Id     uint64
    Type   MessageType
    Target string
    Method string
    Error  string
    Header map[string]string
}复制代码
```

​     Codec接口的Write方法就是编码过程，两个Read是解码过程。

### Registry

​     服务的注册和发现，目前实现的consul,mdns, etcd,etcdv3,zookeeper,kubernetes.等等，

##### Interface
```go
// The registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Register(*Service, ...RegisterOption) error
	Deregister(*Service) error
	GetService(string) ([]*Service, error)
	ListServices() ([]*Service, error)
	Watch() (Watcher, error)
	String() string
}
```

##### Options
```go
type Options struct {
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type RegisterOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
```
​     简单来说，就是Service 进行Register，来进行注册，Client 使用watch方法进行监控，当有服务加入或者删除时这个方法会被触发，以提醒客户端更新Service信息。

​     默认的是服务注册和发现是mdns。

 

### Selector

​    以Registry为基础，Selector 是客户端级别的负载均衡，当有客户端向服务发送请求时， selector根据不同的算法从Registery中的主机列表，得到可用的Service节点，进行通信。目前实现的有循环算法和随机算法，默认的是随机算法。

##### Interface
```go
// Selector builds on the registry as a mechanism to pick nodes
// and mark their status. This allows host pools and other things
// to be built using various algorithms.
type Selector interface {
	Init(opts ...Option) error
	Options() Options
	// Select returns a function which should return the next node
	Select(service string, opts ...SelectOption) (Next, error)
	// Mark sets the success/error against a node
	Mark(service string, node *registry.Node, err error)
	// Reset returns state back to zero for a service
	Reset(service string)
	// Close renders the selector unusable
	Close() error
	// Name of the selector
	String() string
}
```

##### Options
```go
type Options struct {
	Registry registry.Registry
	Strategy Strategy

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type SelectOptions struct {
	Filters  []Filter
	Strategy Strategy

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
```
​     默认的是实现是本地缓存，当前实现的有blacklist,label,named等方式。

###  Broker

​     Broker是消息发布和订阅的接口。很简单的一个例子，因为服务的节点是不固定的，如果有需要修改所有服务行为的需求，可以使服务订阅某个主题，当有信息发布时，所有的监听服务都会收到信息，根据你的需要做相应的行为。

##### Interface
```go
// Broker is an interface used for asynchronous messaging.
// Its an abstraction over various message brokers
// {NATS, RabbitMQ, Kafka, ...}
type Broker interface {
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Init(...Option) error
	Publish(string, *Message, ...PublishOption) error
	Subscribe(string, Handler, ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(Publication) error

type Message struct {
	Header map[string]string
	Body   []byte
}

// Publication is given to a subscription handler for processing
type Publication interface {
	Topic() string
	Message() *Message
	Ack() error
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}
```

##### Options
```go
type Options struct {
	Addrs     []string
	Secure    bool
	Codec     codec.Codec
	TLSConfig *tls.Config
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
```
​     Broker默认的实现方式是http方式，但是这种方式不要在生产环境用。go-plugins里有很多成熟的消息队列实现方式，有kafka、nsq、rabbitmq、redis，等等。

###  Client

​    Client是请求服务的接口，他封装Transport和Codec进行rpc调用，也封装了Brocker进行信息的发布。

##### Interface
```go
// Client is the interface used to make requests to services.
// It supports Request/Response via Transport and Publishing via the Broker.
// It also supports bidiectional streaming of requests.
type Client interface {
	Init(...Option) error
	Options() Options
	NewPublication(topic string, msg interface{}) Publication
	NewRequest(service, method string, req interface{}, reqOpts ...RequestOption) Request
	NewProtoRequest(service, method string, req interface{}, reqOpts ...RequestOption) Request
	NewJsonRequest(service, method string, req interface{}, reqOpts ...RequestOption) Request
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
	CallRemote(ctx context.Context, addr string, req Request, rsp interface{}, opts ...CallOption) error
	Stream(ctx context.Context, req Request, opts ...CallOption) (Streamer, error)
	StreamRemote(ctx context.Context, addr string, req Request, opts ...CallOption) (Streamer, error)
	Publish(ctx context.Context, p Publication, opts ...PublishOption) error
	String() string
}

// Publication is the interface for a message published asynchronously
type Publication interface {
	Topic() string
	Message() interface{}
	ContentType() string
}

// Request is the interface for a synchronous request used by Call or Stream
type Request interface {
	Service() string
	Method() string
	ContentType() string
	Request() interface{}
	// indicates whether the request will be a streaming one rather than unary
	Stream() bool
}

// Streamer is the inteface for a bidirectional synchronous stream
type Streamer interface {
	Context() context.Context
	Request() Request
	Send(interface{}) error
	Recv(interface{}) error
	Error() error
	Close() error
}
```

##### Options
```go
type Options struct {
	// Used to select codec
	ContentType string

	// Plugged interfaces
	Broker    broker.Broker
	Codecs    map[string]codec.NewCodec
	Registry  registry.Registry
	Selector  selector.Selector
	Transport transport.Transport

	// Connection Pool
	PoolSize int
	PoolTTL  time.Duration

	// Middleware for client
	Wrappers []Wrapper

	// Default Call Options
	CallOptions CallOptions

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type CallOptions struct {
	SelectOptions []selector.SelectOption

	// Backoff func
	Backoff BackoffFunc
	// Check if retriable func
	Retry RetryFunc
	// Transport Dial Timeout
	DialTimeout time.Duration
	// Number of Call attempts
	Retries int
	// Request/Response timeout
	RequestTimeout time.Duration

	// Middleware for low level call func
	CallWrappers []CallWrapper

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type RequestOptions struct {
	Stream bool

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
```
​     当然他也支持双工通信 Stream 这些具体的实现方式和使用方式，以后会详细解说。

​     默认的是rpc实现方式，他还有grpc和http方式，在go-plugins里可以找到

### Server

​     Server看名字大家也知道是做什么的了。监听等待rpc请求。监听broker的订阅信息，等待信息队列的推送等。

##### Interface
```go
type Server interface {
	Options() Options
	Init(...Option) error
	Handle(Handler) error
	NewHandler(interface{}, ...HandlerOption) Handler
	NewSubscriber(string, interface{}, ...SubscriberOption) Subscriber
	Subscribe(Subscriber) error
	Register() error
	Deregister() error
	Start() error
	Stop() error
	String() string
}

type Publication interface {
	Topic() string
	Message() interface{}
	ContentType() string
}

type Request interface {
	Service() string
	Method() string
	ContentType() string
	Request() interface{}
	// indicates whether the request will be streamed
	Stream() bool
}

// Streamer represents a stream established with a client.
// A stream can be bidirectional which is indicated by the request.
// The last error will be left in Error().
// EOF indicated end of the stream.
type Streamer interface {
	Context() context.Context
	Request() Request
	Send(interface{}) error
	Recv(interface{}) error
	Error() error
	Close() error
}
```

##### Options
```go
type Options struct {
	Codecs       map[string]codec.NewCodec
	Broker       broker.Broker
	Registry     registry.Registry
	Transport    transport.Transport
	Metadata     map[string]string
	Name         string
	Address      string
	Advertise    string
	Id           string
	Version      string
	HdlrWrappers []HandlerWrapper
	SubWrappers  []SubscriberWrapper

	RegisterTTL time.Duration

	// Debug Handler which can be set by a user
	DebugHandler debug.DebugHandler

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
```
​     默认的是rpc实现方式，他还有grpc和http方式，在go-plugins里可以找到

 

### Service

​     Service是Client和Server的封装，他包含了一系列的方法使用初始值去初始化Service和Client，使我们可以很简单的创建一个rpc服务。

源码：

```
type Service interface {
    Init(...Option)
    Options() Options
    Client() client.Client
    Server() server.Server
    Run() error
    String() string
}
```

### Wrapper

micro 在处理 client 或者 server 的handler的时候会先将 装饰器 Wrapper 执行。

常用的装饰器有：

- 日志
- JWT鉴权
- 熔断，如：hystrix
- metrics，如：prometheus
- 链路追踪，如：jaeger

### 默认值

```
Transport: http
Codec: Protobuf
Registry: mdns
Selector: cache
Broker: http
Client: rpc
Server: rpc
```

