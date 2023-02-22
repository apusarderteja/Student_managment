package handler

import (
	"log"
	"net/http"
)


// user create Template parese
// user create Template parese


func (h Handler) pareseCreateUserTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-user.html")
	if t == nil {
		log.Println("unable to lookup create-user template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// user edit Template parese
// user edit Template parese


func (h Handler) pareseEditUserTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("edit-user.html")
	if t == nil {
		log.Println("unable to lookup edit-user template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

}

// login  Template parese
// login  Template parese


func (h Handler) pareseLoginTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("login.html")
	if t == nil {
		log.Println("unable to lookup login template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// class create Template parese
// class create Template parese


func (h Handler) pareseClassTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create_class.html")
	if t == nil {
		log.Println("unable to lookup class create template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// class Edit Template parese
// class Edit Template parese


func (h Handler) pareseEditClassTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("edit_class.html")
	if t == nil {
		log.Println("unable to lookup class-edit template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

}


// Subject Edit Template parese
// Subject Edit Template parese

func (h Handler) pareseEditSubjectTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("subject_edit.html")
	if t == nil {
		log.Println("unable to lookup Subject-edit template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

}


// Subject create Template parese
// Subject create Template parese

func (h Handler) pareseAddSubjectTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("subject_create.html")
	if t == nil {
		log.Println("unable to lookup add subject template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Student create Template parese
// Student create Template parese
func (h Handler) pareseCreateStudentTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("student_create.html")
	if t == nil {
		log.Println("unable to lookup student create template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Student Edit Template parese
// Student Edit Template parese
func (h Handler) pareseStudentEditTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("student_edit.html")
	if t == nil {
		log.Println("unable to lookup student create template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
func (h Handler) PareseMarkTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("add_mark.html")
	if t == nil {
		log.Println("unable to lookup add mark template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}