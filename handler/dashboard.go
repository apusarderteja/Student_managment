package handler

import (
	"net/http"

	"github.com/justinas/nosurf"
)


func (h Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	h.DashboardTemplate(w, UserForm{
		CSRFToken: nosurf.Token(r),
	})
}
