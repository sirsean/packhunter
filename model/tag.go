package model

import (
	"github.com/sirsean/friendly-ph/ph"
	"gopkg.in/mgo.v2/bson"
)

type Tag struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	Owner  BasicUser
	Name   string
	Public bool
	Users  []ph.User
}

type BasicTag struct {
	Id        string
	OwnerId   string
	Name      string
	UserCount int
}

type BasicTagSubscribed struct {
	BasicTag
	Subscribed bool
}

func (t *Tag) Basic() BasicTag {
	return BasicTag{
		Id:        t.Id.Hex(),
		OwnerId:   t.Owner.Id,
		Name:      t.Name,
		UserCount: len(t.Users),
	}
}

func (t *Tag) AddUser(user ph.User) {
	if !t.HasUser(user) {
		t.Users = append(t.Users, user)
	}
}

func (t *Tag) RemoveUser(user ph.User) {
	if index := t.IndexOfUser(user); index != -1 {
		t.Users = append(t.Users[:index], t.Users[index+1:]...)
	}
}

func (t *Tag) HasUser(user ph.User) bool {
	return t.IndexOfUser(user) != -1
}

func (t *Tag) IndexOfUser(user ph.User) int {
	for i, u := range t.Users {
		if u.Id == user.Id {
			return i
		}
	}
	return -1
}
