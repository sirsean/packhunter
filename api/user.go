package api

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sirsean/friendly-ph/web"
	"github.com/sirsean/friendly-ph/ph"
	"github.com/sirsean/friendly-ph/mongo"
)

func ShowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	session := mongo.Session()
	defer session.Close()

	currentUser, _ := web.CurrentUser(r, session)

	user := ph.GetUserByUsername(currentUser.AccessToken, username)

	response, _ := json.Marshal(user)
	w.Write(response)
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	web.Logout(w, r)
}
