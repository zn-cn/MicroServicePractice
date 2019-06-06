package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	pb "Ethan/MicroServicePractice/interface-center/out/vessel"
	"Ethan/MicroServicePractice/vessel/dao"
)

const (
	DB_NAME           = "MicroServicePractice"
	VESSEL_COLLECTION = "vessels"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
	Create(*pb.Vessel) error
	Close()
}

type VesselRepository struct {
	session *mgo.Session
}

func GetVesselRepository(session *mgo.Session) *VesselRepository {
	return &VesselRepository{session: session}
}

// 接口实现
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	// 选择最近一条容量、载重都符合的货轮
	vessel := dao.Vessel{}
	err := dao.Find(repo.collection(), bson.M{
		"capacity":   bson.M{"$gte": spec.Capacity},
		"max_weight": bson.M{"$gte": spec.MaxWeight},
	}, &vessel)
	if err != nil {
		return nil, err
	}
	pbVessel := Vessel2PBVessel(&vessel)
	return &pbVessel, nil
}

func (repo *VesselRepository) Create(vessel *pb.Vessel) error {
	data := PBVessel2Vessel(vessel)
	return dao.Insert(repo.collection(), &data)
}

func (repo *VesselRepository) Close() {
	repo.session.Close()
}

func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(VESSEL_COLLECTION)
}

func PBVessel2Vessel(vessel *pb.Vessel) dao.Vessel {
	data := dao.Vessel{
		Name:      vessel.Name,
		Capacity:  vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Available: vessel.Available,
		OwerId:    vessel.OwerId,
	}
	if bson.IsObjectIdHex(vessel.Id) {
		data.Id = bson.ObjectIdHex(vessel.Id)
	}
	return data
}

func Vessel2PBVessel(vessel *dao.Vessel) pb.Vessel {
	return pb.Vessel{
		Name:      vessel.Name,
		Capacity:  vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Available: vessel.Available,
		OwerId:    vessel.OwerId,
	}
}
