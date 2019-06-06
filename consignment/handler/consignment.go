package handler

import (
	"context"
	"fmt"
	"log"

	json "github.com/json-iterator/go"

	"github.com/micro/go-micro/broker"
	"gopkg.in/mgo.v2"

	"github.com/yun-mu/MicroServicePractice/common"
	logPB "github.com/yun-mu/MicroServicePractice/interface-center/out/log"

	"github.com/yun-mu/MicroServicePractice/config"
	"github.com/yun-mu/MicroServicePractice/consignment/model"
	pb "github.com/yun-mu/MicroServicePractice/interface-center/out/consignment"
	userPb "github.com/yun-mu/MicroServicePractice/interface-center/out/user"
	vesselPb "github.com/yun-mu/MicroServicePractice/interface-center/out/vessel"
)

// 微服务服务端 struct handler 必须实现 protobuf 中定义的 rpc 方法
// 实现方法的传参等可参考生成的 consignment.pb.go
type Handler struct {
	session      *mgo.Session
	vesselClient vesselPb.VesselServiceClient
	userClient   userPb.UserServiceClient
	Broker       broker.Broker
}

const service = "consignment"
var (
	topic       string
	serviceName string
    version string
)

func init() {
	topic = config.GetBrokerTopic("log")
	serviceName = config.GetServiceName(service)
	version = config.GetVersion(service)
	if version == "" {
		version = "latest"
	}
}

func GetHandler(session *mgo.Session, vesselClient vesselPb.VesselServiceClient, userClient userPb.UserServiceClient, bk broker.Broker) *Handler {
	return &Handler{
		session:      session,
		vesselClient: vesselClient,
		userClient:   userClient,
		Broker:       bk,
	}
}

// 从主会话中 Clone() 出新会话处理查询
func (h *Handler) GetRepo() model.Repository {
	return model.GetConsignmentRepository(h.session.Clone())
}

func (h *Handler) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	repo := h.GetRepo()
	defer repo.Close()
	// 检查是否有适合的货轮
	vReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}

	vResp, err := h.vesselClient.FindAvailable(ctx, vReq, common.Filter(version))
	if err != nil {
		return err
	}

	// 货物被承运
	req.VesselId = vResp.Vessel.Id
	err = repo.Create(req)
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = req

	// 后置操作
	go func() {
		userID := ""
		if id, ok := ctx.Value("user_id").(string); ok {
			userID = id
		}
		msg := fmt.Sprintf("found vessel: %s\n", vResp.Vessel.Name)
		h.pubLog(userID, "CreateConsignment", msg)
	}()
	return nil
}

func (h *Handler) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	repo := h.GetRepo()
	defer repo.Close()
	consignments, err := repo.GetAll()
	if err != nil {
		return err
	}
	resp.Consignments = consignments
	go func() {
		userID := ""
		if id, ok := ctx.Value("user_id").(string); ok {
			userID = id
		}
		msg := fmt.Sprintf("consignments len: %d\n", len(consignments))
		h.pubLog(userID, "GetConsignments", msg)
	}()
	return nil
}

// 发送log
func (h *Handler) pubLog(userID, method, msg string) error {
	logPB := logPB.Log{
		Method: method,
		Origin: serviceName,
		Msg:    msg,
	}
	body, err := json.Marshal(logPB)
	if err != nil {
		return err
	}

	data := &broker.Message{
		Header: map[string]string{
			"user_id": userID,
		},
		Body: body,
	}

	if err := h.Broker.Publish(topic, data); err != nil {
		log.Printf("[pub] failed: %v\n", err)
	}
	return nil
}
