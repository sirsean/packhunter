package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirsean/packhunter/config"
	"github.com/sirsean/packhunter/model"
	"github.com/sirsean/packhunter/mongo"
	"github.com/sirsean/packhunter/ph"
	"github.com/sirsean/packhunter/service"
	"github.com/sirsean/packhunter/web"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var indexTemplate = buildTemplate("index.html")

func Index(w http.ResponseWriter, r *http.Request) {
	session := mongo.Session()
	defer session.Close()

	user, _ := web.CurrentUser(r, session)
	log.Printf("index for %v", user)

	type Data struct {
		UserId   string
		Username string
	}
	data := Data{
		UserId:   user.UserId(),
		Username: user.Me.Username,
	}
	indexTemplate.Execute(w, data)
}

var showUserTemplate = buildTemplate("show-user.html")

func ShowUser(w http.ResponseWriter, r *http.Request) {
	session := mongo.Session()
	defer session.Close()

	user, _ := web.CurrentUser(r, session)

	vars := mux.Vars(r)
	username := vars["username"]

	type Data struct {
		UserId       string
		Username     string
		ShowUsername string
	}
	data := Data{
		UserId:       user.UserId(),
		Username:     user.Me.Username,
		ShowUsername: username,
	}
	showUserTemplate.Execute(w, data)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	log.Printf("/signin")
	url := fmt.Sprintf("%v/oauth/authorize?client_id=%v&redirect_uri=%v&response_type=code&scope=public private",
		config.Get().ProductHunt.Endpoint,
		config.Get().ProductHunt.ApiKey,
		fmt.Sprintf("%v/signin_redirect", config.Get().Host.Name))
	http.Redirect(w, r, url, 302)
}

func SigninRedirect(w http.ResponseWriter, r *http.Request) {
	log.Printf("/signin_redirect")
	log.Printf("%v", r.URL.Query())
	code := r.URL.Query().Get("code")
	log.Printf("code: %v", code)

	payload := map[string]string{
		"client_id":     config.Get().ProductHunt.ApiKey,
		"client_secret": config.Get().ProductHunt.ApiSecret,
		"redirect_uri":  fmt.Sprintf("%v/signin_redirect", config.Get().Host.Name),
		"code":          code,
		"grant_type":    "authorization_code",
	}
	log.Printf("payload: %v", payload)
	jsonPayload, _ := json.Marshal(payload)
	log.Printf("payload: %v", string(jsonPayload))

	url := fmt.Sprintf("%v/oauth/token", config.Get().ProductHunt.Endpoint)
	log.Printf("url: %v", url)
	authReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	authReq.Header.Set("Accept", "application/json")
	authReq.Header.Set("Content-Type", "application/json")
	authReq.Header.Set("Host", "localhost")
	authReq.Header.Set("Cookie", "")

	client := &http.Client{}
	resp, err := client.Do(authReq)
	if err != nil {
		log.Printf("failed %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("body: %v", body)

	b := make(map[string]interface{})
	json.Unmarshal(body, &b)
	log.Printf("b: %v", b)

	accessToken := b["access_token"].(string)
	log.Printf("access token: %v", accessToken)

	phUser := ph.Me(accessToken)
	log.Printf("logged in: %v", phUser)

	session := mongo.Session()
	defer session.Close()

	user, err := service.GetUserByPHId(session, phUser.Id)
	if err != nil {
		user = model.User{
			Id:   bson.NewObjectId(),
			PHId: phUser.Id,
			Me:   phUser,
		}
	}
	user.AccessToken = accessToken
	service.SaveUser(session, &user)

	service.EnsureFollowingTag(session, &user)
	service.SyncFollowingTag(session, &user)

	web.Login(w, r, user)
	http.Redirect(w, r, "/", 302)
}

func buildTemplate(file string) *template.Template {
	return template.Must(template.ParseFiles(templatePath(file)))
}

func templatePath(file string) string {
	return fmt.Sprintf("%s/template/%s", config.Get().Host.Path, file)
}
