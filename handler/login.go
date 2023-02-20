package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Student_managment/Project/storage/postgres"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
	"golang.org/x/crypto/bcrypt"
)

type LoginUser struct {
	Username  string
	Password  string
	FormError map[string]error
	CSRFToken string
}

func (lu LoginUser) Validate() error {
	return validation.ValidateStruct(&lu,
		validation.Field(&lu.Username,
			validation.Required.Error("The username field is required."),
		),
		validation.Field(&lu.Password,
			validation.Required.Error("The password field is required."),
		),
	)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.pareseLoginTemplate(w, LoginUser{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	var lf LoginUser
	if err := h.decoder.Decode(&lf, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := lf.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			formErr := make(map[string]error)
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
			lf.FormError = formErr
			lf.Password = ""
			lf.CSRFToken = nosurf.Token(r)
			h.pareseLoginTemplate(w, lf)
			return
		}
	}

	user, err := h.storage.GetUserByUsername(lf.Username)
	if err != nil {
		if err.Error() == postgres.NotFound {
			formErr := make(map[string]error)
			formErr["Username"] = fmt.Errorf("credentials does not match")
			lf.FormError = formErr
			lf.CSRFToken = nosurf.Token(r)
			lf.Password = ""
			h.pareseLoginTemplate(w, lf)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lf.Password)); err != nil {
		formErr := make(map[string]error)
		formErr["Username"] = fmt.Errorf("credentials does not match")
		lf.FormError = formErr
		lf.CSRFToken = nosurf.Token(r)
		lf.Password = ""
		h.pareseLoginTemplate(w, lf)
		return
	}

	h.sessionManager.Put(r.Context(), "userID", strconv.Itoa(user.ID))
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}


