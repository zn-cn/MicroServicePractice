package dao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Consignment struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`
	Weight      int32         `json:"weight,omitempty" bson:"weight,omitempty"`
	Containers  []*Container  `json:"containers,omitempty" bson:"containers,omitempty"`
	VesselId    string        `json:"vessel_id,omitempty" bson:"vessel_id,omitempty"`
}

type Container struct {
	Id         string `json:"id,omitempty" bson:"id,omitempty"`
	CustomerId string `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	Origin     string `json:"origin,omitempty" bson:"origin,omitempty"`
	UserId     string `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

func Insert(collection *mgo.Collection, consignment *Consignment) error {
	if consignment.Id == "" {
		consignment.Id = bson.NewObjectId()
	}
	return collection.Insert(consignment)
}

func FindAll(collection *mgo.Collection, query interface{}, consignments *[]*Consignment) error {
	return collection.Find(query).All(consignments)
}
