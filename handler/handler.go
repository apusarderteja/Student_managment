package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"Student_managment/Project/storage"
	"github.com/Masterminds/sprig"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/form"
)

type Handler struct {
	sessionManager *scs.SessionManager
	decoder        *form.Decoder
	storage        dbStorage
	Templates      *template.Template
}

type dbStorage interface {
	ListUser(storage.UserFilter) ([]storage.User, error)
	CreateUser(storage.User) (*storage.User, error)
	UpdateUser(storage.User) (*storage.User, error)
	GetUserByID(string) (*storage.User, error)
	GetUserByUsername(string) (*storage.User, error)
	DeleteUserByID(string) error

	DeleteClassByID(id string) error
	CreateClass(storage.Class) (*storage.Class, error)
	ListClass() ([]storage.Class, error)
	GetClassByIDQuery() ([]storage.Class, error)
	GetclassIDByIDQuery(id string) (*storage.Class, error)
	Updateclass(storage.Class) (*storage.Class, error)

	AddSubject(storage.Subject) (*storage.Subject, error)
	ListSubjectQuery() ([]storage.Subject, error)
	GetsubjectIDByIDQuery(id string) (*storage.Subject, error)
	UpdateSubjectFUNC(storage.Subject) (*storage.Subject, error)
	DeleteSubjectByIdQuery(id string) error

	CreateStudent(storage.Student) (*storage.Student, error)
	ListStudentQuery() ([]storage.Student, error)
	DeleteStudentByIdQuery(id string) error
	GetstudentIDByIDQuery(id string) (*storage.Student, error)
	UpdateStudent(u storage.Student) (*storage.Student, error)

	//StudentSubject
	CreateStudentSubject(s storage.StudentSubject) (*storage.StudentSubject, error)
	GetSubIdBYID(classID int) ([]storage.Subject, error)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

type ErrorPage struct {
	Code    int
	Message string
}

func (h Handler) Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	ep := ErrorPage{
		Code:    code,
		Message: error,
	}

	tf := "default"
	switch code {
	case 400, 401, 402, 403, 404:
		tf = "4xx"
	case 500, 501, 503:
		tf = "5xx"
	}

	tpl := fmt.Sprintf("%s.html", tf)
	t := h.Templates.Lookup(tpl)
	if t == nil {
		log.Fatalln("unable to find template")
	}

	if err := t.Execute(w, ep); err != nil {
		log.Fatalln(err)
	}
}

func NewHandler(sm *scs.SessionManager, formDecoder *form.Decoder, storage dbStorage) *chi.Mux {
	h := &Handler{
		sessionManager: sm,
		decoder:        formDecoder,
		storage:        storage,
	}

	h.ParseTemplates()
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(Method)

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Get("/", h.Home)
		r.Get("/login", h.Login)
		r.Post("/login", h.LoginPostHandler)
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets/src"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Use(h.Authentication)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.ListUser)

			r.Get("/create", h.CreateUser)
			r.Post("/store", h.StoreUser)
			r.Get("/{id:[0-9]+}/edit", h.EditUser)
			r.Put("/{id:[0-9]+}/update", h.UpdateUser)
			r.Get("/{id:[0-9]+}/delete", h.DeleteUser)
		})

		r.Route("/class", func(r chi.Router) {
			r.Get("/classlist", h.ListClass)
			r.Get("/{id:[0-9]+}/edit", h.EditClass)
			r.Put("/{id:[0-9]+}/edit", h.UpdateClass)
			r.Get("/{id:[0-9]+}/delete", h.DeleteClass)
			r.Get("/create", h.CreateClass)
			r.Post("/store", h.StoreClass)

		})
		r.Route("/sub", func(r chi.Router) {

			r.Get("/subject", h.AddSubject)
			r.Post("/substore", h.SubjectStore)
			r.Get("/sublist", h.ListSubject)
			r.Get("/{id:[0-9]+}/editsub", h.EditSubject)
			r.Put("/{id:[0-9]+}/editsub", h.UpdateSubject)
			r.Get("/{id:[0-9]+}/delete", h.DeleteSubject)

		})
		r.Route("/student", func(r chi.Router) {

			r.Get("/", h.CreateStudent)
			r.Post("/studentstore", h.StudentStore)
			r.Get("/studentlist", h.ListStudent)
			r.Get("/{id:[0-9]+}/edit-student", h.EditStudent)
			r.Put("/{id:[0-9]+}/edit-student", h.UpdateStudent)
			r.Get("/{id:[0-9]+}/deletestudent", h.DeleteStudent)

		})

		r.Get("/logout", h.LogoutHandler)
	})

	return r
}

func Method(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch strings.ToLower(r.PostFormValue("_method")) {
			case "put":
				r.Method = http.MethodPut
			case "patch":
				r.Method = http.MethodPatch
			case "delete":
				r.Method = http.MethodDelete
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := h.sessionManager.GetString(r.Context(), "userID")
		uID, err := strconv.Atoi(userID)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if uID <= 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) ParseTemplates() error {
	templates := template.New("web-templates").Funcs(template.FuncMap{
		"calculatePreviousPage": func(currentPageNumber int) int {
			if currentPageNumber == 1 {
				return 0
			}

			return currentPageNumber - 1
		},

		"calculateNextPage": func(currentPageNumber, totalPage int) int {
			if currentPageNumber == totalPage {
				return 0
			}

			return currentPageNumber + 1
		},
	}).Funcs(sprig.FuncMap())

	newFS := os.DirFS("assets/templates")
	tmpl := template.Must(templates.ParseFS(newFS, "*/*/*.html", "*/*.html", "*.html"))
	if tmpl == nil {
		log.Fatalln("unable to parse templates")
	}

	h.Templates = tmpl
	return nil
}
