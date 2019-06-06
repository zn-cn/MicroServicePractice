package main

import (
	"log"

	"Ethan/MicroServicePractice/common"
	pb "Ethan/MicroServicePractice/interface-center/out/user"

	"Ethan/MicroServicePractice/user/handler"
)

const service = "user"

func main() {
	session, err := common.CreateDBSession(service)
	if err != nil {
		log.Fatalf("create session error: %v\n", err)
	}
	// 创建于 MongoDB 的主会话，需在退出 main() 时候手动释放连接
	defer session.Close()

	srv := common.GetMicroServer(service)

	bk := srv.Server().Options().Broker
	pb.RegisterUserServiceHandler(srv.Server(), handler.GetHandler(session, bk))

	if err := srv.Run(); err != nil {
		log.Fatalf("user service error: %v\n", err)
	}

}
