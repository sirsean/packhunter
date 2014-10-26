package api

import (
	"net/http"
	"strings"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirsean/friendly-ph/mongo"
	"github.com/sirsean/friendly-ph/model"
	"github.com/sirsean/friendly-ph/web"
	"github.com/sirsean/friendly-ph/service"
	"github.com/sirsean/friendly-ph/rank"
	"github.com/sirsean/friendly-ph/ph"
	"encoding/json"
)

var postDecoder = schema.NewDecoder()

func ListMyTags(w http.ResponseWriter, r *http.Request) {
	session := mongo.Session()
	defer session.Close()

	user, _ := web.CurrentUser(r, session)

	tags := make([]model.BasicTag, len(user.Tags))
	for i, t := range user.Tags {
		tag, _ := service.GetTagByIdHex(session, t.Id)
		tags[i] = tag.Basic()
	}

	response, _ := json.Marshal(tags)
	w.Write(response)
}

func ShowTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session := mongo.Session()
	defer session.Close()

	tag, _ := service.GetTagByIdHex(session, id)

	response, _ := json.Marshal(tag)
	w.Write(response)
}

func CreateTag(w http.ResponseWriter, r *http.Request) {
	type CreateForm struct {
		Name string `schema:"name"`
		Public bool `schema:"public"`
	}

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	r.ParseForm()
	form := new(CreateForm)
	postDecoder.Decode(form, r.PostForm)

	tag := model.Tag{
		Name: form.Name,
		Public: form.Public,
	}

	service.CreateTag(session, &currentUser, &tag)

	response, _ := json.Marshal(tag)
	w.Write(response)
}

func GetTagProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	tag, _ := service.GetTagByIdHex(session, id)

	products := rank.ForTag(currentUser.AccessToken, tag)

	response, _ := json.Marshal(products)
	w.Write(response)
}

func SubscribeToTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	tag, _ := service.GetTagByIdHex(session, id)

	currentUser.AddTag(tag)
	service.SaveUser(session, &currentUser)
}

func UnsubscribeFromTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	tag, _ := service.GetTagByIdHex(session, id)

	currentUser.RemoveTag(tag)
	service.SaveUser(session, &currentUser)
}

func SetTagUsers(w http.ResponseWriter, r *http.Request) {
	type UsersForm struct {
		Usernames string `schema:"usernames"`
	}

	vars := mux.Vars(r)
	id := vars["id"]

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	r.ParseForm()
	form := new(UsersForm)
	postDecoder.Decode(form, r.PostForm)
	usernames := strings.Split(form.Usernames, ",")

	tag, _ := service.GetTagByIdHex(session, id)

	for _, u := range tag.Users {
		if !usernamesContains(usernames, u.Username) {
			tag.RemoveUser(u)
		}
	}
	for _, username := range usernames {
		user := ph.GetUserByUsername(currentUser.AccessToken, username)
		tag.AddUser(user)
	}
	service.SaveTag(session, &tag)
}

func usernamesContains(usernames []string, username string) bool {
	for _, u := range usernames {
		if u == username {
			return true
		}
	}
	return false
}
