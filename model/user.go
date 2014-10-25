package model

import(
	"gopkg.in/mgo.v2/bson"
	"github.com/sirsean/friendly-ph/ph"
	"errors"
)

type User struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	PHId int `bson:"phid"`
	AccessToken string `bson:"access_token"`
	Me ph.User `bson:"me"`
	Tags []BasicTag
}

type BasicUser struct {
	Id string
	Username string
}

func (u *User) UserId() string {
	return u.Id.Hex()
}

func (u *User) Basic() BasicUser {
	return BasicUser{
		Id: u.Id.Hex(),
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
	u.Tags = append(u.Tags, tag.Basic())
}
