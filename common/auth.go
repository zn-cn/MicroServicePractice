package common

import (
	"github.com/yun-mu/MicroServicePractice/config"
	userPb "github.com/yun-mu/MicroServicePractice/interface-center/out/user"
	"context"
	"errors"
	"log"
	"os"

	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
)

var (
	userClient userPb.UserServiceClient
)

func GetUserClient() userPb.UserServiceClient {
	if userClient == nil {
		srv := GetMicroClient("user")
		userClient = userPb.NewUserServiceClient(config.GetServiceName("user"), srv.Client())
	}
	return userClient
}

//
//  AuthWrapper 是一个高阶函数，入参是 "下一步" 函数，出参是认证函数
// 在返回的函数内部处理完认证逻辑后，再手动调用 fn() 进行下一步处理
// token 是从 上下文中取出的，再调用 user-service 将其做验证
// 认证通过则 fn() 继续执行，否则报错
//
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	log.Println("AuthWrapper")
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		// consignment-service 独立测试时不进行认证
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		token := meta["token"]

		// Auth here
		authResp, err := GetUserClient().ValidateToken(context.Background(), &userPb.Token{
			Token: token,
		})

		log.Println("Auth Resp:", authResp)
		if err != nil {
			return err
		}


		// 这里将 JWT 解析出来的 user_id 传递下去
		ctx = context.WithValue(ctx, "user_id", authResp.UserId)
		err = fn(ctx, req, resp)
		return err
	}
}
