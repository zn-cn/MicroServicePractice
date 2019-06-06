package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"

	json "github.com/json-iterator/go"

	"Ethan/MicroServicePractice/common"
	"Ethan/MicroServicePractice/config"
	pb "Ethan/MicroServicePractice/interface-center/out/vessel"

	"github.com/micro/go-micro/metadata"
)

const (
	DEFAULT_INFO_FILE = "vessel.json"
	service           = "vessel"
)

var (
	prefixPath  string
	serviceName string
	token       string
)

func init() {
	serviceName = config.GetServiceName(service)
	prefixPath = os.Getenv("PrefixPath")
	if prefixPath == "" {
		gopath := os.Getenv("GOPATH")
		prefixPath = gopath + "/src/Ethan/MicroServicePractice/"
	}
	token = os.Getenv("Token")
}

// 读取 consignment.json 中记录的货物信息
func parseFile(fileName string) ([]pb.Vessel, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var vessels []pb.Vessel
	err = json.Unmarshal(data, &vessels)
	if err != nil {
		return nil, errors.New("vessel.json file content error")
	}
	return vessels, nil
}

func main() {
	srv := common.GetMicroClient(service)

	// 创建微服务的客户端，简化了手动 Dial 连接服务端的步骤
	client := pb.NewVesselServiceClient(serviceName, srv.Client())

	infoFile := prefixPath + "/vessel-cli/" + DEFAULT_INFO_FILE

	// 解析货物信息
	vessels, err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	// 创建带有用户 token 的 context
	// consignment-service 服务端将从中取出 token，解密取出用户身份
	tokenContext := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	// 调用 RPC
	// 将货物存储到指定用户的仓库里
	for _, vessel := range vessels {
		resp, err := client.Create(tokenContext, &vessel)
		if err != nil {
			log.Fatalf("create vessel error: %v", err)
		}
		log.Printf("created: %t", resp.Created)
	}
}
