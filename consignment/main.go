package main

import (
	"log"

	"Ethan/MicroServicePractice/common"
	"Ethan/MicroServicePractice/config"

	"github.com/micro/go-micro"

	"Ethan/MicroServicePractice/consignment/handler"
	pb "Ethan/MicroServicePractice/interface-center/out/consignment"
	userPb "Ethan/MicroServicePractice/interface-center/out/user"
	vesselPb "Ethan/MicroServicePractice/interface-center/out/vessel"
)

const service = "consignment"

func main() {
	session, err := common.CreateDBSession(service)
	if err != nil {
		log.Fatalf("create session error: %v\n", err)
	}
	// 创建于 MongoDB 的主会话，需在退出 main() 时候手动释放连接
	defer session.Close()

	srv := common.GetMicroServer(service, micro.WrapHandler(common.AuthWrapper))

	// 作为 vessel-service 的客户端
	vClient := vesselPb.NewVesselServiceClient(config.GetServiceName("vessel"), srv.Client())
	uClient := userPb.NewUserServiceClient(config.GetServiceName("user"), srv.Client())
	bk := srv.Server().Options().Broker
	h := handler.GetHandler(session, vClient, uClient, bk)

	// 将 server 作为微服务的服务端
	pb.RegisterShippingServiceHandler(srv.Server(), h)

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
