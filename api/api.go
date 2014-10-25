package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirsean/friendly-ph/mongo"
	"github.com/sirsean/friendly-ph/model"
	"github.com/sirsean/friendly-ph/web"
	"github.com/sirsean/friendly-ph/service"
	"github.com/sirsean/friendly-ph/rank"
	"encoding/json"
)

var postDecoder = schema.NewDecoder()

func ListMyTags(w http.ResponseWriter, r *http.Request) {
	session := mongo.Session()
	defer session.Close()

	user, _ := web.CurrentUser(r, session)

	response, _ := json.Marshal(user.Tags)
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
