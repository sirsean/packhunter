package model

import (
	"errors"
	"github.com/sirsean/friendly-ph/ph"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	PHId        int           `bson:"phid"`
	AccessToken string        `bson:"access_token"`
	Me          ph.User       `bson:"me"`
	Tags        []BasicTag
}

type BasicUser struct {
	Id       string
	Username string
}

func (u *User) UserId() string {
	return u.Id.Hex()
}

func (u *User) Basic() BasicUser {
	return BasicUser{
		Id:       u.Id.Hex(),
		Username: u.Me.Username,
	}
}

func (u *User) Tag(name string) (BasicTag, error) {
	for _, t := range u.Tags {
		if t.Name == name {
			return t, nil
		}
	}
	err := errors.New("tag not found")
	return BasicTag{}, err
}

func (u *User) AddTag(tag Tag) {
	if !u.HasTag(tag.Basic()) {
		u.Tags = append(u.Tags, tag.Basic())
	}
}

func (u *User) RemoveTag(tag Tag) {
	if index := u.IndexOfTag(tag.Basic()); index != -1 {
		u.Tags = append(u.Tags[:index], u.Tags[index+1:]...)
	}
}

func (u *User) HasTag(tag BasicTag) bool {
	return u.IndexOfTag(tag) != -1
}

func (u *User) IndexOfTag(tag BasicTag) int {
	for i, t := range u.Tags {
		if t.Id == tag.Id {
			return i
		}
	}
	return -1
}
