package storage

import (
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserFilter struct {
	SearchTerm string
	Offset     int
	Limit      int
}

type User struct {
	ID        int          `json:"id" form:"-" db:"id"`
	FirstName string       `json:"first_name" db:"first_name"`
	LastName  string       `json:"last_name" db:"last_name"`
	Email     string       `json:"email" db:"email"`
	Username  string       `json:"username" db:"username"`
	Password  string       `json:"password" db:"password"`
	Status    bool         `json:"status" db:"status"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
	Total     int          `json:"-" db:"total"`
	FormError map[string]error
}

type Subject struct {
	ID          int          `json:"id" form:"-" db:"id"`
	ClassID     int          `json:"class_id" db:"class_id"`
	ClassName   string       `json:"class_name" db:"class_name"`
	SubjectName string       `json:"subject_name" db:"subject_name"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at" db:"deleted_at"`
	FormError   map[string]error
}

type Class struct {
	ID        int          `json:"id" form:"-" db:"id"`
	ClassName string       `json:"class_name" db:"class_name"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
	FormError map[string]error
}
type Student struct {
	ID        int          `json:"id" form:"-" db:"id"`
	ClassID   int          `json:"class_id" db:"class_id"`
	ClassName string       `json:"class_name" db:"class_name"`
	FirstName string       `json:"first_name" db:"first_name"`
	LastName  string       `json:"last_name" db:"last_name"`
	Roll      int          `json:"roll" db:"roll"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
	FormError map[string]error
}

type StudentSubject struct {
	ID          int          `json:"id" form:"-" db:"id"`
	StudentId   int          `json:"student_id" db:"student_id"`
	SubjectId   int          `json:"subject_id" db:"subject_id"`
	ClassName string       `json:"class_name" db:"class_name"`
	FirstName string       `json:"first_name" db:"first_name"`
	LastName  string       `json:"last_name" db:"last_name"`
	Roll      int          `json:"roll" db:"roll"`
	Marks       int          `json:"marks" db:"marks"`
	SubjectName string       `json:"subject_name" db:"subject_name"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at" db:"deleted_at"`
	Marksa        map[int]int
	FormError map[string]error
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName,
			validation.Required.Error("The first name field is required."),
			validation.Length(3, 20).Error("The first name length must be between 3 and 20 characters."),
		),
		validation.Field(&u.LastName,
			validation.Required.Error("The last name field is required."),
			validation.Length(3, 20).Error("The last name length must be between 3 and 20 characters."),
		),
		validation.Field(&u.Username,
			validation.Required.When(u.ID == 0).Error("The username field is required."),
			validation.Length(3, 20).Error("The User name  length must be between 3 and 20 characters."),
		),
		validation.Field(&u.Email,
			validation.Required.When(u.ID == 0).Error("The email field is required."),
			is.Email.Error("The email field must be a valid email."),
		),
		validation.Field(&u.Password,
			validation.Required.When(u.ID == 0).Error("The password field is required."),
		),
	)

}

func (c Class) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClassName,
			validation.Required.Error("The class name field is required."),
		),
	)
}

func (s Subject) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.SubjectName,
			validation.Required.Error("The subject name field is required."),
		),
	)
}
func (s Student) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.FirstName,
			validation.Required.Error("The First name field is required."),
			validation.Length(3, 20).Error("The first name length must be between 3 and 20 characters."),
		),
		validation.Field(&s.LastName,
			validation.Required.Error("The Last name field is required."),
			validation.Length(3, 20).Error("The last name length must be between 3 and 20 characters."),
		),
		validation.Field(&s.Roll,
			validation.Required.Error("The Roll field is required."),
			validation.Min(1).Error("The Roll must be a positive number."),
			validation.Max(200).Error("only 200 rolls are allowed."),
		),
	)
}
