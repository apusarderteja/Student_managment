package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"


	"Student_managment/Project/storage"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"

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
	
	class := storage.Class{}
	if err := h.decoder.Decode(&class, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

if err := class.Validate(); err != nil {
	if vErr, ok := err.(validation.Errors); ok {
		class.FormError = vErr
		fmt.Println(class.FormError)
	}
	h.pareseClassTemplate(w, ClassForm{
		Class:     class,
		CSRFToken: nosurf.Token(r),
		FormError: class.FormError,
	})
	return
}

checkAlreadyExist, err := h.IsClassAlreadyExistsCheck(w, r, class.ClassName)

	if err != nil {
		fmt.Println(err)
		return
	}
	if checkAlreadyExist {
		h.pareseClassTemplate(w, ClassForm{
			Class:     class,
			CSRFToken: nosurf.Token(r),
			FormError: map[string]error{
				"class_name": fmt.Errorf("The Class already Exist."),
			}})
		return
	}

	_, eXrr := h.storage.CreateClass(class)
	if eXrr != nil {
		log.Println(eXrr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/class/classlist", http.StatusSeeOther)

}


func (h Handler) ListClass(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	
	listCl, err := h.storage.ListClass()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("class_list.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

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


// For Class Already Exists Check By CLassname
func (h Handler) IsClassAlreadyExistsCheck(w http.ResponseWriter, r *http.Request, ClassName string) (bool, error) {
	Class, err := h.storage.CheckClassExists(ClassName)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return Class, nil
}
