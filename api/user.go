package api

import (
	"strings"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sirsean/friendly-ph/web"
	"github.com/sirsean/friendly-ph/ph"
	"github.com/sirsean/friendly-ph/mongo"
	"github.com/sirsean/friendly-ph/service"
	"github.com/sirsean/friendly-ph/model"
)

func ShowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	user := ph.GetUserByUsername(currentUser.AccessToken, username)

	tags, _ := service.TagsUserIsOn(session, currentUser, user)

	type UserResponse struct {
		ph.User
		Tags []model.BasicTag `json:"tags"`
	}

	response, _ := json.Marshal(UserResponse{
		User: user,
		Tags: tags,
	})
	w.Write(response)
}

func SetUserTags(w http.ResponseWriter, r *http.Request) {
	type TagsForm struct {
		TagIds string `schema:"tag_ids"`
	}

	vars := mux.Vars(r)
	username := vars["username"]

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	user := ph.GetUserByUsername(currentUser.AccessToken, username)

	r.ParseForm()
	form := new(TagsForm)
	postDecoder.Decode(form, r.PostForm)
	tagIds := strings.Split(form.TagIds, ",")

	for _, t := range currentUser.Tags {
		tag, _ := service.GetTagByIdHex(session, t.Id)
		if tagsContains(tagIds, t.Id) {
			tag.AddUser(user)
		} else {
			tag.RemoveUser(user)
		}
		service.SaveTag(session, &tag)
	}
}

func tagsContains(tagIds []string, tagId string) bool {
	for _, t := range tagIds {
		if t == tagId {
			return true
		}
	}
	return false
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	web.Logout(w, r)
}
