package service

import (
	"github.com/sirsean/packhunter/config"
	"github.com/sirsean/packhunter/model"
	"github.com/sirsean/packhunter/ph"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func SaveTag(session *mgo.Session, tag *model.Tag) {
	coll := tagCollection(session)
	coll.UpdateId(tag.Id, tag)
}

func TagsUserIsOn(session *mgo.Session, currentUser model.User, user ph.User) ([]model.BasicTag, error) {
	tags := make([]model.BasicTag, 0)
	for _, t := range currentUser.Tags {
		tag, err := GetTagByIdHex(session, t.Id)
		if err != nil {
			return nil, err
		}
		for _, u := range tag.Users {
			if u.Username == user.Username {
				tags = append(tags, t)
			}
		}
	}
	return tags, nil
}

func UserPublicTags(session *mgo.Session, currentUser model.User, user ph.User) ([]model.BasicTagSubscribed, error) {
	tags := make([]model.BasicTagSubscribed, 0)
	mUser, err := GetUserByPHId(session, user.Id)
	if err != nil {
		return nil, err
	}
	for _, t := range mUser.Tags {
		tag, err := GetTagByIdHex(session, t.Id)
		if err != nil {
			return nil, err
		}

		// if this user is the owner and the tag is public
		if tag.Owner.Id == mUser.Id.Hex() && tag.Public {
			tags = append(tags, model.BasicTagSubscribed{
				BasicTag:   t,
				Subscribed: currentUser.HasTag(t),
			})
		}
	}
	return tags, nil
}

func EnsureFollowingTag(session *mgo.Session, user *model.User) {
	if _, err := user.Tag("Following"); err != nil {
		tag := model.Tag{
			Name:   "Following",
			Public: false,
		}
		CreateTag(session, user, &tag)
	}
}

func SyncFollowingTag(session *mgo.Session, user *model.User) {
	coll := tagCollection(session)

	if basicTag, err := user.Tag("Following"); err == nil {
		tag, _ := GetTagByIdHex(session, basicTag.Id)

		users := ph.Following(user.AccessToken, user.PHId)
		tag.Users = users

		coll.UpdateId(tag.Id, tag)
	}
}
