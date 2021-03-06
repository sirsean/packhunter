package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirsean/packhunter/api"
	"github.com/sirsean/packhunter/config"
	"github.com/sirsean/packhunter/controller"
	"github.com/sirsean/packhunter/mongo"
	"github.com/sirsean/packhunter/service"
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
	router.HandleFunc("/signin", controller.Signin).Methods("GET")
	router.HandleFunc("/signin_redirect", controller.SigninRedirect).Methods("GET")
	router.HandleFunc("/user/{username}", controller.ShowUser).Methods("GET")

	router.HandleFunc("/api/users", api.ListMyUsers).Methods("GET")
	router.HandleFunc("/api/tags/mine", api.ListMyTags).Methods("GET")
	router.HandleFunc("/api/tags/{id}", api.ShowTag).Methods("GET")
	router.HandleFunc("/api/tags/{id}/products", api.GetTagProducts).Methods("GET")
	router.HandleFunc("/api/tags/{id}/subscribe", api.SubscribeToTag).Methods("POST")
	router.HandleFunc("/api/tags/{id}/unsubscribe", api.UnsubscribeFromTag).Methods("POST")
	router.HandleFunc("/api/tags/{id}/users", api.SetTagUsers).Methods("PUT")
	router.HandleFunc("/api/tags", api.CreateTag).Methods("POST")
	router.HandleFunc("/api/user/{username}", api.ShowUser).Methods("GET")
	router.HandleFunc("/api/user/{username}/tags", api.SetUserTags).Methods("PUT")

	router.HandleFunc("/api/logout", api.UserLogout).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(fmt.Sprintf("%s/static/", config.Get().Host.Path))))
	http.Handle("/", router)

	port := config.Get().Host.Port
	log.Printf("Serving on port %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
