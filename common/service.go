package common

// default setting:
//   Transport: tcp
//   Server: rpc
//   Client: rpc
//   RegisterTTL: 30s
//   RegisterInterval: 20s
//   Registry: consul
//   Broker: kafka
//   Selector: cache
//   Codec: protobuf
//   Tracing: jaeger
//   Metrics: jaeger
//   breaker: hystrix 注：客户端熔断
//   ratelimit: uber/ratelimit

import (
	"github.com/yun-mu/MicroServicePractice/config"
	"log"
	"os"
	"time"

	gh "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/broker/kafka"
	"github.com/micro/go-plugins/transport/tcp"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
)

var (
	defaultOpts    []micro.Option
	defaultWebOpts []web.Option
	defaultServer  micro.Option
	defaultClient  micro.Option
)

func init() {
	defaultOpts = []micro.Option{
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 20),
		micro.Transport(tcp.NewTransport()),
	}
	defaultWebOpts = []web.Option{
		web.RegisterTTL(time.Second * 30),
		web.RegisterInterval(time.Second * 20),
	}
	gh.DefaultMaxConcurrent = 100
	gh.DefaultVolumeThreshold = 50
}

func GetMicroClient(service string, exOpts ...micro.Option) micro.Service {
	opts := getOpts(service, exOpts...)
	if defaultClient != nil {
		opts = append(opts, defaultClient)
	}
	t, _, err := NewJaegerTracer(config.GetServiceName(service), config.GetTracingAddr(service))
	if err != nil {
		log.Fatalf("opentracing tracer create error:%v", err)
	}
	opts = append(opts, micro.WrapClient(hystrix.NewClientWrapper(), opentracing.NewClientWrapper(t), ratelimit.NewClientWrapper(1024)))
	srv := micro.NewService(opts...)
	// 解析命令行参数
	srv.Init()
	return srv
}

func GetMicroServer(service string, exOpts ...micro.Option) micro.Service {
	opts := getOpts(service, exOpts...)
	if defaultServer != nil {
		opts = append(opts, defaultServer)
	}
	brokerKafka := kafka.NewBroker(func(options *broker.Options) {
		options.Addrs = config.GetBrokerAddrs(service)
	})
	if err := brokerKafka.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
	t, _, err := NewJaegerTracer(config.GetServiceName(service), config.GetTracingAddr(service))
	if err != nil {
		log.Fatalf("opentracing tracer create error:%v", err)
	}
	opts = append(opts, micro.Broker(brokerKafka), micro.WrapHandler(opentracing.NewHandlerWrapper(t), ratelimit.NewHandlerWrapper(1024)))
	srv := micro.NewService(opts...)
	// 解析命令行参数
	srv.Init()
	return srv
}

func getOpts(service string, exOpts ...micro.Option) []micro.Option {
	opts := append(exOpts,
		defaultOpts...,
	)
	version := config.GetVersion(service)
	if version == "" {
		version = "latest"
	}

	serviceName := config.GetServiceName(service)
	opts = append(opts, micro.Version(version), micro.Name(serviceName))
	if os.Getenv("DebugMDNS") == "" {
		opts = append(opts, micro.Registry(consul.NewRegistry(func(op *registry.Options) {
			op.Addrs = config.GetRegistryAddrs(service)
		})))
	}
	return opts
}

func GetMicroWeb(service string, exOpts ...web.Option) web.Service {
	opts := append(exOpts,
		defaultWebOpts...,
	)
	version := config.GetVersion(service)
	if version == "" {
		version = "latest"
	}
	serviceName := config.GetServiceName(service)
	opts = append(opts, web.Version(version), web.Name(serviceName))
	if os.Getenv("DebugMDNS") == "" {
		opts = append(opts, web.Registry(consul.NewRegistry(func(op *registry.Options) {
			op.Addrs = config.GetRegistryAddrs(service)
		})))
	}

	srv := web.NewService(opts...)
	// 解析命令行参数
	srv.Init()
	return srv
}
