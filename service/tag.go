package service

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sirsean/friendly-ph/config"
	"github.com/sirsean/friendly-ph/model"
	"github.com/sirsean/friendly-ph/ph"
)

var tagCollection = func(session *mgo.Session) *mgo.Collection {
	return session.DB(config.Get().Mongo.Database).C("tags")
}

func GetTagByIdHex(session *mgo.Session, id string) (model.Tag, error) {
	coll := tagCollection(session)
	var tag model.Tag
	err := coll.FindId(bson.ObjectIdHex(id)).One(&tag)
	return tag, err
}

func CreateTag(session *mgo.Session, user *model.User, tag *model.Tag) {
	coll := tagCollection(session)

	tag.Id = bson.NewObjectId()
	tag.Owner = user.Basic()
	coll.Insert(tag)

	user.AddTag(*tag)
	SaveUser(session, user)
}

func EnsureUncategorizedTag(session *mgo.Session, user *model.User) {
	if _, err := user.Tag("#Uncategorized"); err != nil {
		tag := model.Tag{
			Name: "#Uncategorized",
			Public: false,
		}
		CreateTag(session, user, &tag)
	}
}

func SyncUncategorizedTag(session *mgo.Session, user *model.User) {
	coll := tagCollection(session)

	if basicTag, err := user.Tag("#Uncategorized"); err == nil {
		tag, _ := GetTagByIdHex(session, basicTag.Id)

		users := ph.Following(user.AccessToken, user.PHId)
		tag.Users = users

		coll.UpdateId(tag.Id, tag)
	}
}
