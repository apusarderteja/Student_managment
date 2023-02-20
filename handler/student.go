package handler

import (
	// "html/template"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "strings"

	"Student_managment/Project/storage"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

type StudentForm struct {
	Student    storage.Student
	Classlists []storage.Class
	Class      storage.Class
	FormError  map[string]error
	CSRFToken  string
}

func (h Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	classlist, err := h.storage.GetClassByIDQuery()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseCreateStudentTemplate(w, StudentForm{
		Classlists: classlist,
		CSRFToken:  nosurf.Token(r),
	})
}

func (h Handler) StudentStore(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	student := storage.Student{}

	// cl := storage.Class{}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}
	fmt.Println(student)

	classlist, err := h.storage.GetClassByIDQuery()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr

		}
		h.pareseCreateStudentTemplate(w, StudentForm{
			Student:    student,
			Classlists: classlist,
			CSRFToken:  nosurf.Token(r),
			FormError:  student.FormError,
		})
		return
	}

	stID, erra := h.storage.CreateStudent(student)
	if erra != nil {
		log.Println(erra)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println(stID.ID)
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

	erraa := h.StuSub(w, r, student.ClassID, stID.ID)
	if erraa != nil {
		log.Println(erraa)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/student/studentlist", http.StatusSeeOther)

}



func (h Handler) ListStudent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	liststudent, err := h.storage.ListStudentQuery()
	fmt.Println(liststudent)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("student_list.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, liststudent); err != nil {
		log.Println(err)
	}
}

func (h Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.storage.DeleteStudentByIdQuery(id); err != nil {
		h.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/student/studentlist", http.StatusSeeOther)
}

func (h Handler) EditStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	studen, err := h.storage.GetstudentIDByIDQuery(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseStudentEditTemplate(w, StudentForm{
		Student:   *studen,
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	// var form StudentForm
	student := storage.Student{ID: uID}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	fmt.Println("#################################")
	fmt.Println(student)
	fmt.Println("#################################")
	// form.Student = student
	// if err := student.Validate(); err != nil {
	// 	if vErr, ok := err.(validation.Errors); ok {
	// 		student.FormError = vErr
	// 		fmt.Println(student.FormError)
	// 	}
	// 	h.pareseEditSubjectTemplate(w, SubjectForm{
	// 		Student:     student,
	// 		CSRFToken: nosurf.Token(r),
	// 		FormError: student.FormError,
	// 	})
	// 	return
	// }

	_, err = h.storage.UpdateStudent(student)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/student/studentlist", http.StatusSeeOther)
}

func (h Handler) StuSub(w http.ResponseWriter, r *http.Request, classID int, stuID int) error {
	sub, _ := h.storage.GetSubIdBYID(classID)
	for _, val := range sub {
		dbphrm := storage.StudentSubject{
			StudentId: stuID,
			SubjectId: val.ID,
			Marks:     0,
		}
		_, err := h.storage.CreateStudentSubject(dbphrm)
		fmt.Println(err)

	}
	return nil

}