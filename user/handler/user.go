package handler

import (
	"context"
	"errors"
	"fmt"
	"log"

	json "github.com/json-iterator/go"

	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"

	"github.com/yun-mu/MicroServicePractice/config"
	logPB "github.com/yun-mu/MicroServicePractice/interface-center/out/log"
	pb "github.com/yun-mu/MicroServicePractice/interface-center/out/user"

	"github.com/yun-mu/MicroServicePractice/user/model"
)

type Handler struct {
	session      *mgo.Session
	tokenService model.Authable
	Broker       broker.Broker
}

var (
	topic       string
	serviceName string
)

func init() {
	topic = config.GetBrokerTopic("log")
	serviceName = config.GetServiceName("user")
}

func GetHandler(session *mgo.Session, bk broker.Broker) *Handler {
	return &Handler{
		session:      session,
		tokenService: model.GetTokenService(),
		Broker:       bk,
	}
}

// 从主会话中 Clone() 出新会话处理查询
func (h *Handler) GetRepo() model.Repository {
	return model.GetUserRepository(h.session.Clone())
}

func (h *Handler) Create(ctx context.Context, req *pb.User, resp *pb.Response) error {
	// 哈希处理用户输入的密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPwd)

	repo := h.GetRepo()
	defer repo.Close()
	if err := repo.Create(req); err != nil {
		return nil
	}
	resp.User = req

	go func() {
		msg := fmt.Sprintf("[user] id:%s name: %s email: %s Created", req.GetId(), req.GetName(), req.GetEmail())
		h.pubLog(req.GetId(), "Create", msg)
	}()

	return nil
}

func (h *Handler) Get(ctx context.Context, req *pb.User, resp *pb.Response) error {
	repo := h.GetRepo()
	defer repo.Close()
	u, err := repo.Get(req.GetId())
	if err != nil {
		return err
	}
	resp.User = u

	go func() {
		msg := fmt.Sprintf("[user] id:%s", req.GetId())
		h.pubLog(req.GetId(), "Get", msg)
	}()

	return nil
}

func (h *Handler) GetAll(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	repo := h.GetRepo()
	defer repo.Close()
	users, err := repo.GetAll()
	if err != nil {
		return err
	}
	resp.Users = users

	go func() {
		msg := ""
		h.pubLog("", "GetAll", msg)
	}()

	return nil
}

func (h *Handler) Auth(ctx context.Context, req *pb.User, resp *pb.Token) error {
	// 在 part3 中直接传参 &pb.User 去查找用户
	// 会导致 req 的值完全是数据库中的记录值
	// 即 req.Password 与 u.Password 都是加密后的密码
	// 将无法通过验证
	repo := h.GetRepo()
	defer repo.Close()
	u, err := repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	// 进行密码验证
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return err
	}
	t, err := h.tokenService.Encode(u)
	if err != nil {
		return err
	}
	resp.Token = t

	go func() {
		msg := fmt.Sprintf("[user] email: %s", req.GetEmail())
		h.pubLog("", "Auth", msg)
	}()

	return nil
}

func (h *Handler) ValidateToken(ctx context.Context, req *pb.Token, resp *pb.Token) error {
	// Decode token
	claims, err := h.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}
	if claims.User.Id == "" {
		return errors.New("invalid user")
	}
	resp.Valid = true
	resp.UserId = claims.User.Id

	go func() {
		msg := fmt.Sprintf("[user] id: %s", claims.User.Id)
		h.pubLog(claims.User.Id, "ValidateToken", msg)
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
		log.Fatalf("[pub] failed: %v\n", err)
	}
	return nil
}
