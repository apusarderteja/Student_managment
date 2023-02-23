package handler

import (
	"log"
	"net/http"
)
//For LogOut User
func (h Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionManager.Destroy(r.Context()); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
