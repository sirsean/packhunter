package model

import(
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
}

func (u *User) UserId() string {
	return u.Id.Hex()
}
