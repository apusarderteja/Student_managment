package handler

import (
	"Student_managment/Project/storage"
	
	"log"
	"net/http"


	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
)

type MarkForm struct {
	Mark      []storage.StudentSubject
	StudentId string
	FormError map[string]error
	CSRFToken string
}

func (h Handler) AddMark(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	mark, err := h.storage.GetStudentIdBySubjectID(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("add_mark.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	data := MarkForm{
		Mark: mark,
		// FormError: map[string]error{},
		CSRFToken: nosurf.Token(r),
		StudentId: id,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func (h Handler) MarkStore(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	markk := storage.StudentSubject{}

	if err := h.decoder.Decode(&markk, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}


	for t, y := range markk.Marksa {
		dbphrm := storage.StudentSubject{
			StudentId: markk.StudentId,
			SubjectId: t,
			Marks:     y,
		}

		_, err := h.storage.UpdateStudentMark(dbphrm)
		if err != nil {
			log.Println(err)
			// http.Error(w, "internal server error", http.StatusInternalServerError)
		}

	}

	http.Redirect(w, r, "/student/studentlist", http.StatusSeeOther)

}
