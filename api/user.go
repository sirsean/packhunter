package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirsean/packhunter/model"
	"github.com/sirsean/packhunter/mongo"
	"github.com/sirsean/packhunter/ph"
	"github.com/sirsean/packhunter/service"
	"github.com/sirsean/packhunter/web"
	"net/http"
	"strings"
)

func ListMyUsers(w http.ResponseWriter, r *http.Request) {
	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	t, _ := currentUser.Tag("Following")
	tag, _ := service.GetTagByIdHex(session, t.Id)

	response, _ := json.Marshal(tag.Users)
	w.Write(response)
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	user := ph.GetUserByUsername(currentUser.AccessToken, username)

	tags, _ := service.TagsUserIsOn(session, currentUser, user)
	publicTags, _ := service.UserPublicTags(session, currentUser, user)

	type UserResponse struct {
		ph.User
		Tags       []model.BasicTag           `json:"tags"`
		PublicTags []model.BasicTagSubscribed `json:"public_tags"`
	}

	response, _ := json.Marshal(UserResponse{
		User:       user,
		Tags:       tags,
		PublicTags: publicTags,
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
