package api

import (
	"net/http"
	"github.com/sirsean/friendly-ph/web"
)

func UserLogout(w http.ResponseWriter, r *http.Request) {
	web.Logout(w, r)
}
