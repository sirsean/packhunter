package mongo

import (
	"github.com/sirsean/packhunter/config"
	"gopkg.in/mgo.v2"
	"log"
	//"gopkg.in/mgo.v2/bson"
	"time"
)

var session *mgo.Session

func Connect() {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{config.Get().Mongo.Hosts},
		Timeout:  60 * time.Second,
		Database: config.Get().Mongo.AuthDatabase,
		Username: config.Get().Mongo.AuthUsername,
		Password: config.Get().Mongo.AuthPassword,
	}
	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
		return
	}
	session.SetMode(mgo.Monotonic, true)
}

func Session() *mgo.Session {
	return session.Copy()
}
