package dao

import (
	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name,omitempty" bson:"name,omitempty"`
	Company  string        `json:"company,omitempty" bson:"company,omitempty"`
	Email    string        `json:"email,omitempty" bson:"email,omitempty"`
	Password string        `json:"password,omitempty" bson:"password,omitempty"`
}

func Insert(collection *mgo.Collection, user *User) error {
	if user.Id == "" {
		user.Id = bson.NewObjectId()
	}
	return collection.Insert(user)
}

func FindAll(collection *mgo.Collection, query interface{}, users *[]*User) error {
	return collection.Find(query).All(users)
}

func Find(collection *mgo.Collection, query interface{}, user *User) error {
	return collection.Find(query).One(user)
}
