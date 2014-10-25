package web

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/sirsean/friendly-ph/config"
	"github.com/sirsean/friendly-ph/model"
	//"github.com/sirsean/friendly-ph/service"
	"gopkg.in/mgo.v2"
	"io"
	"log"
	"net/http"
	"strings"
)

var cookieStore *sessions.CookieStore

func init() {
	cookieStore = sessions.NewCookieStore([]byte(config.Get().CookieStore.AuthenticationKey))
	cookieStore.Options = &sessions.Options{
		MaxAge:   86400 * 365,
		HttpOnly: true,
		Path:     "/",
	}
}

func Login(w http.ResponseWriter, r *http.Request, user model.User) {
	id := newSessionID()
	log.Printf("logging in user %v with session %v", user.Id.Hex(), id)
	session, _ := getSessionByName(r, id)
	session.Values["UserId"] = user.Id.Hex()
	if err := session.Save(r, w); err != nil {
		log.Printf("failed to save session: %v", err)
	}
	cookie := sessions.NewCookie("session_id", id, cookieStore.Options)
	http.SetCookie(w, cookie)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getCurrentSessionID(r)
	log.Printf("getting session %v", sessionId)
	if err != nil {
		return
	}
	session, err := getSessionByName(r, sessionId)
	if err != nil {
		return
	}
	session.Values["UserId"] = nil
	if err := session.Save(r, w); err != nil {
		log.Printf("failed to save session: %v", err)
	}
}

func newSessionID() string {
	c := 128
	b := make([]byte, c)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		panic("Could not generate random number")
	}
	return strings.Replace(base64.URLEncoding.EncodeToString(b), "=", "", -1)
}

func getSessionByName(r *http.Request, name string) (session *sessions.Session, err error) {
	session, err = cookieStore.Get(r, name)
	return
}

func getCurrentSessionID(r *http.Request) (id string, err error) {
	var c *http.Cookie
	if c, err = r.Cookie("session_id"); err != nil {
		return
	}
	id = c.Value
	return
}

func CurrentUser(r *http.Request, s *mgo.Session) (user model.User, err error) {
	sessionId, err := getCurrentSessionID(r)
	log.Printf("getting session %v", sessionId)
	if err != nil {
		return
	}
	session, err := getSessionByName(r, sessionId)
	if err != nil {
		return
	}
	userId := session.Values["UserId"]
	if userId == nil {
		err = errors.New("no user")
		return
	}
	err = errors.New("no user!")
	//user, err = service.GetUserByIdHex(s, userId.(string))
	return
}
