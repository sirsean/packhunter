package service

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sirsean/friendly-ph/config"
	"github.com/sirsean/friendly-ph/model"
)

var userCollection = func(session *mgo.Session) *mgo.Collection {
	return session.DB(config.Get().Mongo.Database).C("users")
}

func GetUserByIdHex(session *mgo.Session, id string) (model.User, error) {
	coll := userCollection(session)
	var user model.User
	err := coll.FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func GetUserByPHId(session *mgo.Session, phUserId int) (model.User, error) {
	coll := userCollection(session)
	var user model.User
	err := coll.Find(bson.M{"phid": phUserId}).One(&user)
	return user, err
}

func SaveUser(session *mgo.Session, user *model.User) error {
	coll := userCollection(session)
	_, err := coll.Upsert(bson.M{"phid": user.PHId}, user)
	return err
}
