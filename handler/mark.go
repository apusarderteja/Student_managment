package handler

import (
	"Student_managment/Project/storage"
	"fmt"



	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
)

type MarkForm struct {
	Mark      []storage.StudentSubject
	ProfileWithMark      []storage.StudentSubject
	ProfiledataWithMark      storage.StudentSubject
	StudentId string
	FormError map[string]error
	CSRFToken string
}
// For Add student mark
// For Add student mark
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
		CSRFToken: nosurf.Token(r),
		StudentId: id,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}
// For  student mark STORE
// For  student mark STORE

func (h Handler) MarkStore(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	studentdetails := storage.StudentSubject{}

	if err := h.decoder.Decode(&studentdetails, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}


	for t, y := range studentdetails.Marksa {
		studetails := storage.StudentSubject{
			StudentId: studentdetails.StudentId,
			SubjectId: t,
			Marks:     y,
		}
		

		_, err := h.storage.UpdateStudentMark(studetails)
		if err != nil {
			log.Println(err)
		}

	}

	http.Redirect(w, r, "/student/studentlist", http.StatusSeeOther)

}

func (h Handler) StuSub(w http.ResponseWriter, r *http.Request, classID int, stuID int) error {
	sub, _ := h.storage.GetSubIdBYID(classID)
	for _, val := range sub {
		studetails := storage.StudentSubject{
			StudentId: stuID,
			SubjectId: val.ID,
			Marks:     0,
		}
		_, err := h.storage.InsertstudentMarkQuery(studetails)
		fmt.Println(err)

	}
	return nil

}

// FOR SHOWING with Mark STUDENT LIST
// FOR SHOWING Mark with STUDENT LIST
func (h Handler) MarkListStudent(w http.ResponseWriter, r *http.Request) {
	
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	liststudentWithMark, err := h.storage.GetStudentSubjectByStudentID()
	if err != nil {
		log.Println(err)
		
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	
	t := h.Templates.Lookup("mark_show.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, liststudentWithMark); err != nil {
		log.Println(err)
	}
}
// STUDENT DELETE SHOW PROCESS
// STUDENT DELETE SHOW PROCESS


func (h Handler) DeleteProfileT(w http.ResponseWriter, r *http.Request ) {
	id := chi.URLParam(r, "id")

	if err := h.storage.DeleteStudentSubjectByID(id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/student/marklist", http.StatusSeeOther)
}


// STUDENT DETAILS SHOW PROCESS
// STUDENT DETAILS SHOW PROCESS
func (h Handler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	id := chi.URLParam(r, "id")


var profileWithMark []storage.StudentSubject
profileWithMark, err := h.storage.GetStudentProfileQuery(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("profile.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	
	data := MarkForm{
		ProfileWithMark:     profileWithMark,
		ProfiledataWithMark: profileWithMark[0],
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}
