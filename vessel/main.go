package main

import (
	"log"

	"github.com/yun-mu/MicroServicePractice/common"
	pb "github.com/yun-mu/MicroServicePractice/interface-center/out/vessel"

	"github.com/yun-mu/MicroServicePractice/vessel/handler"

	"github.com/micro/go-micro"
)

const service = "vessel"

func main() {
	session, err := common.CreateDBSession(service)
	if err != nil {
		log.Fatalf("create session error: %v\n", err)
	}
	// 创建于 MongoDB 的主会话，需在退出 main() 时候手动释放连接
	defer session.Close()

	srv := common.GetMicroServer(service, micro.WrapHandler(common.AuthWrapper))

	bk := srv.Server().Options().Broker
	// 将实现服务端的 API 注册到服务端
	pb.RegisterVesselServiceHandler(srv.Server(), handler.GetHandler(session, bk))

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
