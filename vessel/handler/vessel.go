package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/micro/go-micro/broker"
	"gopkg.in/mgo.v2"

	"github.com/yun-mu/MicroServicePractice/config"
	logPB "github.com/yun-mu/MicroServicePractice/interface-center/out/log"
	pb "github.com/yun-mu/MicroServicePractice/interface-center/out/vessel"

	"github.com/yun-mu/MicroServicePractice/vessel/model"
)

// 实现微服务的服务端
type Handler struct {
	session *mgo.Session
	Broker  broker.Broker
}

var (
	topic       string
	serviceName string
)

func init() {
	topic = config.GetBrokerTopic("log")
	serviceName = config.GetServiceName("vessel")
}

func GetHandler(session *mgo.Session, bk broker.Broker) *Handler {
	return &Handler{
		session: session,
		Broker:  bk,
	}
}

func (h *Handler) GetRepo() model.Repository {
	return model.GetVesselRepository(h.session.Clone())
}

func (h *Handler) FindAvailable(ctx context.Context, req *pb.Specification, resp *pb.Response) error {
	defer h.GetRepo().Close()
	v, err := h.GetRepo().FindAvailable(req)
	if err != nil {
		return err
	}
	resp.Vessel = v

	go func() {
		msg := fmt.Sprintf("found vessel: %s\n", v.Name)
		h.pubLog("", "FindAvailable", msg)
	}()
	return nil
}

func (h *Handler) Create(ctx context.Context, req *pb.Vessel, resp *pb.Response) error {
	defer h.GetRepo().Close()
	if err := h.GetRepo().Create(req); err != nil {
		return err
	}
	resp.Vessel = req
	resp.Created = true
	go func() {
		msg := fmt.Sprintf("vessel: %s\n", req.Name)
		h.pubLog("", "Create", msg)
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
