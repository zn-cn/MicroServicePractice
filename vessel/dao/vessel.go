package dao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Vessel struct {
	Id        bson.ObjectId `json:"id,omitempty"  bson:"_id,omitempty"`
	Capacity  int32         `json:"capacity,omitempty" bson:"capacity,omitempty"`
	MaxWeight int32         `json:"max_weight,omitempty" bson:"max_weight,omitempty"`
	Name      string        `json:"name,omitempty" bson:"name,omitempty"`
	Available bool          `json:"available,omitempty" bson:"available,omitempty"`
	OwerId    string        `json:"ower_id,omitempty" bson:"ower_id,omitempty"`
}

func Insert(collection *mgo.Collection, vessel *Vessel) error {
	if vessel.Id == "" {
		vessel.Id = bson.NewObjectId()
	}
	return collection.Insert(vessel)
}

func FindAll(collection *mgo.Collection, query interface{}, vessels *[]*Vessel) error {
	return collection.Find(query).All(vessels)
}

func Find(collection *mgo.Collection, query interface{}, vessel *Vessel) error {
	return collection.Find(query).One(vessel)
}
