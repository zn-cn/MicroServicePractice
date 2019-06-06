package model

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	pb "Ethan/MicroServicePractice/interface-center/out/user"
	"Ethan/MicroServicePractice/user/dao"
)

type Repository interface {
	Get(id string) (*pb.User, error)
	GetAll() ([]*pb.User, error)
	Create(*pb.User) error
	GetByEmail(email string) (*pb.User, error)
	Close()
}

type UserRepository struct {
	session *mgo.Session
}

const (
	DB_NAME        = "MicroServicePractice"
	CON_COLLECTION = "users"
)

func GetUserRepository(session *mgo.Session) *UserRepository {
	return &UserRepository{session: session}
}

func (repo *UserRepository) Get(id string) (*pb.User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("ID format is wrong")
	}

	query := bson.M{
		"_id": bson.ObjectIdHex((id)),
	}
	user := &dao.User{}
	err := dao.Find(repo.collection(), query, user)
	if err != nil {
		return nil, err
	}
	pbUser := User2PBUser(user)
	return &pbUser, nil
}

func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	users := []*dao.User{}

	err := dao.FindAll(repo.collection(), nil, &users)
	if err != nil {
		return nil, err
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUser := User2PBUser(user)
		pbUsers[i] = &pbUser
	}
	return pbUsers, nil
}

func (repo *UserRepository) Create(u *pb.User) error {
	query := bson.M{
		"email": u.Email,
	}
	user := dao.User{}
	err := dao.Find(repo.collection(), query, &user)
	if err == nil {
		return errors.New("exist")
	}
	data := PBUser2User(u)
	return dao.Insert(repo.collection(), &data)
}

func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	user := &dao.User{}
	err := dao.Find(repo.collection(), bson.M{
		"email": email,
	}, user)
	if err != nil {
		return nil, err
	}
	pbUser := User2PBUser(user)
	return &pbUser, nil
}

// 关闭连接
func (repo *UserRepository) Close() {
	// Close() 会在每次查询结束的时候关闭会话
	// Mgo 会在启动的时候生成一个 "主" 会话
	// 你可以使用 Copy() 直接从主会话复制出新会话来执行，即每个查询都会有自己的数据库会话
	// 同时每个会话都有自己连接到数据库的 socket 及错误处理，这么做既安全又高效
	// 如果只使用一个连接到数据库的主 socket 来执行查询，那很多请求处理都会阻塞
	// Mgo 因此能在不使用锁的情况下完美处理并发请求
	// 不过弊端就是，每次查询结束之后，必须确保数据库会话要手动 Close
	// 否则将建立过多无用的连接，白白浪费数据库资源
	repo.session.Close()
}

// 返回所有货物信息
func (repo *UserRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(CON_COLLECTION)
}

func PBUser2User(u *pb.User) dao.User {
	data := dao.User{
		Name:     u.Name,
		Company:  u.Company,
		Email:    u.Email,
		Password: u.Password,
	}
	if bson.IsObjectIdHex(u.Id) {
		data.Id = bson.ObjectIdHex(u.Id)
	}
	return data
}

func User2PBUser(u *dao.User) pb.User {
	return pb.User{
		Id:       u.Id.Hex(),
		Name:     u.Name,
		Company:  u.Company,
		Email:    u.Email,
		Password: u.Password,
	}
}
