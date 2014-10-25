package controller

import (
	"fmt"
	//"github.com/gorilla/mux"
	"github.com/sirsean/friendly-ph/config"
	"github.com/sirsean/friendly-ph/model"
	"github.com/sirsean/friendly-ph/mongo"
	//"github.com/sirsean/friendly-ph/service"
	"github.com/sirsean/friendly-ph/web"
	"html/template"
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
		UserId string
		User   model.User
	}
	data := Data{
		UserId: user.UserId(),
		User:   user,
	}
	indexTemplate.Execute(w, data)
}

func buildTemplate(file string) *template.Template {
	return template.Must(template.ParseFiles(templatePath(file)))
}

func templatePath(file string) string {
	return fmt.Sprintf("%s/template/%s", config.Get().Host.Path, file)
}
