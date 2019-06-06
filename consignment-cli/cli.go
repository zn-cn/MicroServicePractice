package main

import (
	"context"

	json "github.com/json-iterator/go"

	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/micro/go-micro/metadata"

	"Ethan/MicroServicePractice/common"
	"Ethan/MicroServicePractice/config"
	pb "Ethan/MicroServicePractice/interface-center/out/consignment"
)

const (
	DEFAULT_INFO_FILE = "consignment.json"
	service           = "consignment"
)

var (
	prefixPath  string
	token       string
	serviceName string
)

func init() {
	prefixPath = os.Getenv("PrefixPath")
	if prefixPath == "" {
		gopath := os.Getenv("GOPATH")
		prefixPath = gopath + "/src/Ethan/MicroServicePractice/"
	}
	token = os.Getenv("Token")
	serviceName = config.GetServiceName(service)
}

// 读取 consignment.json 中记录的货物信息
func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main() {
	srv := common.GetMicroClient(service)

	// 创建微服务的客户端，简化了手动 Dial 连接服务端的步骤
	client := pb.NewShippingServiceClient(serviceName, srv.Client())

	infoFile := prefixPath + "/consignment-cli/" + DEFAULT_INFO_FILE

	// 解析货物信息
	consignment, err := parseFile(infoFile)
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
	resp, err := client.CreateConsignment(tokenContext, consignment)
	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}
	log.Printf("created: %t", resp.Created)

	// 列出目前所有托运的货物
	resp, err = client.GetConsignments(tokenContext, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments: %v", err)
	}
	for i, c := range resp.Consignments {
		log.Printf("consignment_%d: %v\n", i, c)
	}
}
