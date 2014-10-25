package main

import (
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/sirsean/friendly-ph/api"
	"github.com/sirsean/friendly-ph/config"
	"github.com/sirsean/friendly-ph/controller"
	"github.com/sirsean/friendly-ph/mongo"
	"github.com/sirsean/friendly-ph/service"
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting Up")

	mongo.Connect()

	session := mongo.Session()
	service.EnsureIndexes(session)
	session.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", controller.Index).Methods("GET")
	//router.HandleFunc("/doc/{id}", controller.ShowDocument).Methods("GET")

	//router.HandleFunc("/api/docs", api.ListDocuments).Methods("GET")
	//router.HandleFunc("/api/doc", api.CreateDocument).Methods("POST")
	//router.HandleFunc("/api/doc/{id}", api.ShowDocument).Methods("GET")
	//router.HandleFunc("/api/doc/{id}", api.UpdateDocument).Methods("PUT")
	//router.HandleFunc("/api/doc/{id}", api.DeleteDocument).Methods("DELETE")
	//router.HandleFunc("/api/doc/{id}/{section}/{file}", api.DownloadFile).Methods("GET")
	//router.HandleFunc("/api/doc/{id}/{section}/{file}", api.UploadFile).Methods("POST")
	//router.HandleFunc("/api/doc/{id}/compile", api.Compile).Methods("GET")
	//router.HandleFunc("/api/user/signup", api.UserSignup).Methods("POST")
	//router.HandleFunc("/api/user/login", api.UserLogin).Methods("POST")
	//router.HandleFunc("/api/user/logout", api.UserLogout).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(fmt.Sprintf("%s/static/", config.Get().Host.Path))))
	http.Handle("/", router)

	port := config.Get().Host.Port
	log.Printf("Serving on port %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
