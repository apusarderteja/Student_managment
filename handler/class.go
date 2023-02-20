package handler

import (
	// "html/template"
	// "fmt"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "strings"

	"Student_managment/Project/storage"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	// validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

type ClassForm struct {
	Class     storage.Class
	FormError map[string]error
	CSRFToken string
}

func (h Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	h.pareseClassTemplate(w, ClassForm{
		CSRFToken: nosurf.Token(r),
	})
}


func (h Handler) StoreClass(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	
	cl := storage.Class{}
	if err := h.decoder.Decode(&cl, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

if err := cl.Validate(); err != nil {
	if vErr, ok := err.(validation.Errors); ok {
		cl.FormError = vErr
		fmt.Println(cl.FormError)
	}
	h.pareseClassTemplate(w, ClassForm{
		Class:     cl,
		CSRFToken: nosurf.Token(r),
		FormError: cl.FormError,
	})
	return
}

	_, err := h.storage.CreateClass(cl)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/class/classlist", http.StatusSeeOther)

}


func (h Handler) ListClass(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	// var err error
	// CurrentPage := 1
	// pn := r.FormValue("page")
	// if pn != "" {
	// 	CurrentPage, err = strconv.Atoi(pn)
	// 	if err != nil {
	// 		CurrentPage = 1
	// 	}
	// }

	// offset := 0
	// if CurrentPage > 1 {
	// 	offset = (CurrentPage * userListLimit) - userListLimit
	// }

	// st := r.FormValue("SearchTerm")
	// uf := storage.UserFilter{
	// 	SearchTerm: st,
	// 	Offset:     offset,
	// 	Limit:      userListLimit,
	// }
	
	listCl, err := h.storage.ListClass()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("class_list.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	// total := 0
	// if len(listUser) > 0 {
	// 	total = listUser[0].Total
	// }

	// totalPage := int(math.Ceil(float64(total) / float64(userListLimit)))

	// data := UserList{
	// 	Users:       listUser,
	// 	// SearchTerm:  st,
	// 	// CurrentPage: CurrentPage,
	// 	// Limit:       userListLimit,
	// 	// Total:       total,
	// 	// TotalPage:   totalPage,
	// }

	if err := t.Execute(w, listCl); err != nil {
		log.Println(err)
	}
}


func (h Handler) EditClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cl, err := h.storage.GetclassIDByIDQuery(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseEditClassTemplate(w, ClassForm{
		Class:     *cl,
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) UpdateClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	
	var form ClassForm
	class := storage.Class{ID: uID}
	if err := h.decoder.Decode(&class, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}



	form.Class = class
	if err := class.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			class.FormError = vErr
			fmt.Println(class.FormError)
		}
		h.pareseEditClassTemplate(w, ClassForm{
			Class:     class,
			CSRFToken: nosurf.Token(r),
			FormError: class.FormError,
		})
		return
	}

	_, err = h.storage.Updateclass(class)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/class/classlist", http.StatusSeeOther)
}


func (h Handler) DeleteClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.storage.DeleteClassByID(id); err != nil {
		h.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/class/classlist", http.StatusSeeOther)
}
