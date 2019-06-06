package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"Ethan/MicroServicePractice/consignment/dao"
	pb "Ethan/MicroServicePractice/interface-center/out/consignment"
)

const (
	DB_NAME        = "MicroServicePractice"
	CON_COLLECTION = "consignments"
)

type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

type ConsignmentRepository struct {
	session *mgo.Session
}

func GetConsignmentRepository(session *mgo.Session) *ConsignmentRepository {
	return &ConsignmentRepository{session: session}
}

// 接口实现
func (repo *ConsignmentRepository) Create(con *pb.Consignment) error {
	data := PBConsignment2Consignment(con)
	return dao.Insert(repo.collection(), &data)
}

// 获取全部数据
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	cons := []*dao.Consignment{}
	err := dao.FindAll(repo.collection(), nil, &cons)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Consignment, len(cons))
	for i, con := range cons {
		pbCon := Consignment2PBConsignment(con)
		res[i] = &pbCon
	}
	return res, err
}

// 关闭连接
func (repo *ConsignmentRepository) Close() {
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
func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(CON_COLLECTION)
}

func PBConsignment2Consignment(con *pb.Consignment) dao.Consignment {
	containers := make([]*dao.Container, len(con.Containers))
	for i, container := range con.Containers {
		containers[i] = &dao.Container{
			Id:         container.Id,
			CustomerId: container.CustomerId,
			Origin:     container.Origin,
			UserId:     container.UserId,
		}
	}
	data := dao.Consignment{
		VesselId:    con.VesselId,
		Weight:      con.Weight,
		Containers:  containers,
		Description: con.Description,
	}
	if bson.IsObjectIdHex(con.Id) {
		data.Id = bson.ObjectIdHex(con.Id)
	}
	return data
}

func Consignment2PBConsignment(con *dao.Consignment) pb.Consignment {
	containers := make([]*pb.Container, len(con.Containers))
	for j, container := range con.Containers {
		containers[j] = &pb.Container{
			Id:         container.Id,
			CustomerId: container.CustomerId,
			Origin:     container.Origin,
			UserId:     container.UserId,
		}
	}
	return pb.Consignment{
		Id:          con.Id.Hex(),
		VesselId:    con.VesselId,
		Weight:      con.Weight,
		Containers:  containers,
		Description: con.Description,
	}
}
