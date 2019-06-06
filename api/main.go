package main

import (
	"github.com/yun-mu/MicroServicePractice/api/handler"
	"github.com/yun-mu/MicroServicePractice/common"
	"github.com/yun-mu/MicroServicePractice/config"

	"github.com/yun-mu/MicroServicePractice/api/middleware"
	userPb "github.com/yun-mu/MicroServicePractice/interface-center/out/user"

	"log"

	"github.com/gin-gonic/gin"
)

const (
	service     = "api"
	userService = "user"
)

var (
	serviceName     string
	userServiceName string
)

func init() {
	serviceName = config.GetServiceName(service)
	userServiceName = config.GetServiceName(userService)
}

func main() {
	srv := common.GetMicroWeb(service)

	srvC := common.GetMicroClient(userService)
	// 创建 user-service 微服务的客户端
	userClient := userPb.NewUserServiceClient(userServiceName, srvC.Client())

	userHandler := handler.GetUserHandler(userClient)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.Use(middleware.Logger())
	user := v1.Group("/user")
	user.POST("/login", userHandler.Login)
	user.POST("/sign", userHandler.Sign)

	srv.Handle("/", router)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
