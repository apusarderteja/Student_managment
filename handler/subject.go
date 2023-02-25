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

type SubjectForm struct {
	Subject storage.Subject
	Classlists []storage.Class
	Class      storage.Class
	FormError  map[string]error
	CSRFToken  string
}

func (h Handler) AddSubject(w http.ResponseWriter, r *http.Request) {
	classlist, err := h.storage.GetClassByIDQuery()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseAddSubjectTemplate(w, SubjectForm{
		Classlists: classlist,
		CSRFToken:  nosurf.Token(r),
	})
}

func (h Handler) SubjectStore(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	subject := storage.Subject{}
	if err := h.decoder.Decode(&subject, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}

	classlist, err := h.storage.GetClassByIDQuery()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := subject.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			subject.FormError = vErr

		}
		h.pareseAddSubjectTemplate(w, SubjectForm{
			Subject:    subject,
			Classlists: classlist,
			CSRFToken:  nosurf.Token(r),
			FormError:  subject.FormError,
		})
		return
	}


	_, eRr := h.storage.AddSubject(subject)
	if eRr != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/sub/sublist", http.StatusSeeOther)

}

func (h Handler) EditSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sub, err := h.storage.GetsubjectIDByIDQuery(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseEditSubjectTemplate(w, SubjectForm{
		Subject:   *sub,
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	var form SubjectForm
	subject := storage.Subject{ID: uID}
	if err := h.decoder.Decode(&subject, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form.Subject = subject
	if err := subject.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			subject.FormError = vErr
			fmt.Println(subject.FormError)
		}
		h.pareseEditSubjectTemplate(w, SubjectForm{
			Subject:   subject,
			CSRFToken: nosurf.Token(r),
			FormError: subject.FormError,
		})
		return
	}

	_, err = h.storage.UpdateSubjectFUNC(subject)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/sub/sublist", http.StatusSeeOther)
}

func (h Handler) ListSubject(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	listsub, err := h.storage.ListSubjectQuery()
	fmt.Println(listsub)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("subject_list.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, listsub); err != nil {
		log.Println(err)
	}
}

func (h Handler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.storage.DeleteSubjectByIdQuery(id); err != nil {
		h.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sub/sublist", http.StatusSeeOther)
}
