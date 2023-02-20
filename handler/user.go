package handler

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	// "strings"

	"Student_managment/Project/storage"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

const userListLimit = 2

type UserList struct {
	Users       []storage.User
	SearchTerm  string
	CurrentPage int
	Limit       int
	Total       int
	TotalPage   int
}

type UserForm struct {
	User      storage.User
	FormError map[string]error
	CSRFToken string
}

func (h Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	var err error
	CurrentPage := 1
	pn := r.FormValue("page")
	if pn != "" {
		CurrentPage, err = strconv.Atoi(pn)
		if err != nil {
			CurrentPage = 1
		}
	}

	offset := 0
	if CurrentPage > 1 {
		offset = (CurrentPage * userListLimit) - userListLimit
	}

	st := r.FormValue("SearchTerm")
	uf := storage.UserFilter{
		SearchTerm: st,
		Offset:     offset,
		Limit:      userListLimit,
	}
	
	listUser, err := h.storage.ListUser(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("list-user.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	total := 0
	if len(listUser) > 0 {
		total = listUser[0].Total
	}

	totalPage := int(math.Ceil(float64(total) / float64(userListLimit)))

	data := UserList{
		Users:       listUser,
		SearchTerm:  st,
		CurrentPage: CurrentPage,
		Limit:       userListLimit,
		Total:       total,
		TotalPage:   totalPage,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}


func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.pareseCreateUserTemplate(w, UserForm{
		CSRFToken: nosurf.Token(r),
	})
}


func (h Handler) StoreUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := UserForm{}
	user := storage.User{}
	if err := h.decoder.Decode(&user, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	
	form.User = user
	if err := user.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			user.FormError = vErr
			fmt.Println(user.FormError)
		}
		h.pareseCreateUserTemplate(w, UserForm{
			User:     user,
			CSRFToken: nosurf.Token(r),
			FormError: user.FormError,
		})
		return
	}

	newUser, err := h.storage.CreateUser(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%v/edit", newUser.ID), http.StatusSeeOther)
}



func (h Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	editUser, err := h.storage.GetUserByID(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	var form UserForm
	form.User = *editUser
	form.CSRFToken = nosurf.Token(r)
	h.pareseEditUserTemplate(w, form)
}

func (h Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	var form UserForm
	user := storage.User{ID: uID}
	if err := h.decoder.Decode(&user, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}


	form.User = user
	if err := user.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			user.FormError = vErr
			fmt.Println(user.FormError)
		}
		h.pareseEditUserTemplate(w, UserForm{
			User:     user,
			CSRFToken: nosurf.Token(r),
			FormError: user.FormError,
		})
		return
	}

	_, err = h.storage.UpdateUser(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}



func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.storage.DeleteUserByID(id); err != nil {
		h.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
