package model

import(
	"gopkg.in/mgo.v2/bson"
	"github.com/sirsean/friendly-ph/ph"
)

type Tag struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Owner BasicUser
	Name string
	Public bool
	Users []ph.User
}

type BasicTag struct {
	Id string
	Name string
}

func (t *Tag) Basic() BasicTag {
	return BasicTag{
		Id: t.Id.Hex(),
		Name: t.Name,
	}
}
